package usecases

import (
	"fmt"
	"math"
	"notifier/internal/entities"
	"time"
)

const (
	highThreshold    = 10.0
	lowThreshold     = 4.5
	noDateDifference = 10 * time.Minute
	highText         = "Внимание! Сахар выше 10! Текущее значение %.1f ммоль/л"
	lowText          = "Внимание! Сахар ниже 4.5! Текущее значение %.1f ммоль/л"
	noDataText       = "Нет новых данных"
	bolusText        = "На болюс %.1f введенный в %s значение сахара через %s часа (%s) %.1f ммоль/л"
)

var (
	treatments          map[string]bool
	treatmentTimers     map[string]*time.Timer
	nextTreatmentTimers map[string]*time.Timer
)

type Fetcher interface {
	FetchEntry() (*entities.SGVResponse, error)
	FetchTreatments() (*entities.TreatmentResponse, error)
}

type Notifier interface {
	Send(message string) error
}

type Checker struct {
	fetcher  Fetcher
	notifier Notifier
	lastSent time.Time
	prevSGV  float64
	interval time.Duration
}

func NewChecker(
	fetcher Fetcher,
	notifier Notifier,
	interval time.Duration,
) *Checker {
	treatmentTimers = make(map[string]*time.Timer)
	nextTreatmentTimers = make(map[string]*time.Timer)
	treatments = make(map[string]bool)

	return &Checker{
		fetcher:  fetcher,
		notifier: notifier,
		interval: interval,
	}
}

func (c *Checker) CheckAndNotify() error {
	response, err := c.fetcher.FetchEntry()
	if err != nil {
		return err
	}

	sgv := calculateValue(response.SGV)
	fmt.Println("sgv:", sgv)

	if isLongTimeAgo(response.DateString) {
		if c.isInterval() {
			if err = c.notifier.Send(noDataText); err != nil {
				return err
			}
		}
		c.prevSGV = sgv
		return nil
	}

	if err = c.CheckTreatments(sgv); err != nil {
		return err
	}

	message := c.getMessage(sgv)

	if message != "" {
		if c.isInterval() || (c.prevSGV <= highThreshold && sgv > highThreshold) || (c.prevSGV >= lowThreshold && sgv < lowThreshold) {
			if err = c.notifier.Send(message); err != nil {
				return err
			}
		}
	}

	c.prevSGV = sgv
	return nil
}

func (c *Checker) isInterval() bool {
	if time.Since(c.lastSent) >= c.interval {
		c.lastSent = time.Now()
		return true
	}
	return false
}

func isLongTimeAgo(date time.Time) bool {
	localTime := time.Now().UTC()

	if localTime.Sub(date) > noDateDifference {
		return true
	}
	return false
}

func (c *Checker) CheckTreatments(sgv float64) error {
	response, err := c.fetcher.FetchTreatments()
	if err != nil {
		return err
	}

	if isLongTimeAgo(response.CreatedAt) && response.Insulin == 0 {
		return nil
	}

	if _, ok := treatments[response.UID]; ok {
		return nil
	}

	treatments[response.UID] = true

	if t, exists := treatmentTimers[response.UID]; exists {
		t.Stop()
	}

	fmt.Println("Setting treatment timer for 2 hours")
	treatmentTimers[response.UID] = time.AfterFunc(2*time.Hour, func() {
		formattedTime := response.CreatedAt.Format("15:04")
		nextPeriod := response.CreatedAt.Add(2 * time.Hour)
		formattedNextPeriod := nextPeriod.Format("15:04")
		message := fmt.Sprintf(bolusText, response.Insulin, formattedTime, "2", formattedNextPeriod, sgv)
		c.notifier.Send(message)
		fmt.Println("Timer 2 hours executed")
	})

	if t, exists := nextTreatmentTimers[response.UID]; exists {
		t.Stop()
	}

	fmt.Println("Setting next treatment timer for 4 hours")
	nextTreatmentTimers[response.UID] = time.AfterFunc(4*time.Hour, func() {
		formattedTime := response.CreatedAt.Format("15:04")
		nextPeriod := response.CreatedAt.Add(4 * time.Hour)
		formattedNextPeriod := nextPeriod.Format("15:04")
		message := fmt.Sprintf(bolusText, response.Insulin, formattedTime, "4", formattedNextPeriod, sgv)
		c.notifier.Send(message)
		fmt.Println("Timer 4 hours executed")
	})

	return nil
}

func (c *Checker) getMessage(value float64) string {
	switch {
	case value > highThreshold:
		return fmt.Sprintf(highText, value)
	case value < lowThreshold:
		return fmt.Sprintf(lowText, value)
	case (c.prevSGV > lowThreshold && c.prevSGV < highThreshold) && (value < lowThreshold || value > highThreshold):
		if value < lowThreshold {
			return fmt.Sprintf(lowText, value)
		}
		return fmt.Sprintf(highText, value)
	default:
		return ""
	}
}

func calculateValue(sgv int) float64 {
	return math.Round(float64(sgv)/18*10) / 10
}
