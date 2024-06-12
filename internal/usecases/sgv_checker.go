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
	treatmentsDate   = 20 * time.Minute
	highText         = "Внимание! Сахар выше 10! Текущее значение %.1f ммоль/л"
	lowText          = "Внимание! Сахар ниже 4.5! Текущее значение %.1f ммоль/л"
	noDataText       = "Нет новых данных"
	bolusText        = "На болюс %.1f введенный в %s значение сахара через %s часа (%s) %.1f ммоль/л"
	utcPeriod        = 3 * time.Hour
	firstPeriod      = 2 * time.Hour
	secondPeriod     = 4 * time.Hour
)

var (
	treatments          map[string]bool
	treatmentTimers     map[string]*time.Timer
	nextTreatmentTimers map[string]*time.Timer
)

type Fetcher interface {
	FetchEntry() (*entities.SGVResponse, error)
	FetchTreatments() ([]entities.TreatmentResponse, error)
}

type Notifier interface {
	Send(message string) error
}

type Checker struct {
	fetcher    Fetcher
	notifier   Notifier
	lastSent   time.Time
	prevSGV    float64
	currentSGV float64
	interval   time.Duration
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

	c.currentSGV = calculateValue(response.SGV)
	fmt.Println("sgv:", c.currentSGV)

	if isLongTimeAgo(response.DateString, noDateDifference) {
		if c.isInterval() {
			if err = c.notifier.Send(noDataText); err != nil {
				return err
			}
		}
		c.prevSGV = c.currentSGV
		return nil
	}

	if err = c.CheckTreatments(); err != nil {
		return err
	}

	message := c.getMessage(c.currentSGV)

	if message != "" {
		if c.isInterval() || (c.prevSGV <= highThreshold && c.currentSGV > highThreshold) || (c.prevSGV >= lowThreshold && c.currentSGV < lowThreshold) {
			if err = c.notifier.Send(message); err != nil {
				return err
			}
		}
	}

	c.prevSGV = c.currentSGV
	return nil
}

func (c *Checker) isInterval() bool {
	if time.Since(c.lastSent) >= c.interval {
		c.lastSent = time.Now()
		return true
	}
	return false
}

func isLongTimeAgo(date time.Time, difference time.Duration) bool {
	localTime := time.Now().UTC()

	if localTime.Sub(date) > difference {
		return true
	}
	return false
}

func (c *Checker) CheckTreatments() error {
	responses, err := c.fetcher.FetchTreatments()
	if err != nil {
		return err
	}

	for _, response := range responses {
		if isLongTimeAgo(response.CreatedAt, treatmentsDate) || response.Insulin == 0 {
			continue
		}

		if _, ok := treatments[response.UID]; ok {
			continue
		}

		treatments[response.UID] = true

		if t, exists := treatmentTimers[response.UID]; exists {
			t.Stop()
		}

		fmt.Println("Setting treatment timer for 2 hours")
		treatmentTimers[response.UID] = time.AfterFunc(firstPeriod, func() {
			formattedTime := response.CreatedAt.Format("15:04")
			nextPeriod := response.CreatedAt.Add(utcPeriod).Add(firstPeriod)
			formattedNextPeriod := nextPeriod.Format("15:04")
			message := fmt.Sprintf(bolusText, response.Insulin, formattedTime, "2", formattedNextPeriod, c.currentSGV)
			c.notifier.Send(message)
			fmt.Println("Timer 2 hours executed")
		})

		if t, exists := nextTreatmentTimers[response.UID]; exists {
			t.Stop()
		}

		fmt.Println("Setting next treatment timer for 4 hours")
		nextTreatmentTimers[response.UID] = time.AfterFunc(secondPeriod, func() {
			formattedTime := response.CreatedAt.Format("15:04")
			nextPeriod := response.CreatedAt.Add(utcPeriod).Add(secondPeriod)
			formattedNextPeriod := nextPeriod.Format("15:04")
			message := fmt.Sprintf(bolusText, response.Insulin, formattedTime, "4", formattedNextPeriod, c.currentSGV)
			c.notifier.Send(message)
			fmt.Println("Timer 4 hours executed")
		})
	}

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
