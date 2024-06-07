package usecases

import (
	"fmt"
	"github.com/google/uuid"
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
	bolusText        = "На болюс %.1f введенный в %s значение сахара через 2 часа (%s) %.1f ммоль/л"
)

var (
	treatments          map[uuid.UUID]bool
	treatmentTimers     map[uuid.UUID]*time.Timer
	nextTreatmentTimers map[uuid.UUID]*time.Timer
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

	if isLongTimeAgo(response.DateString) {
		if err = c.notifier.Send(noDataText); err != nil {
			return err
		}
	}

	sgv := calculateValue(response.SGV)
	fmt.Println("sgv:", sgv)

	go c.CheckTreatments(sgv)

	message := c.getMessage(sgv)

	if message != "" {
		if time.Since(c.lastSent) >= c.interval || (c.prevSGV <= highThreshold && sgv > highThreshold) || (c.prevSGV >= lowThreshold && sgv < lowThreshold) {
			if err = c.notifier.Send(message); err != nil {
				return err
			}
			c.lastSent = time.Now()
		}
	}

	c.prevSGV = sgv
	return nil
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

	if !isLongTimeAgo(response.CreatedAt) && response.Insulin == 0 {
		return nil
	}

	treatments[response.UID] = true

	if t, exists := treatmentTimers[response.UID]; exists {
		t.Stop()
	}
	fmt.Println("treatment timer 2 hours is running")
	treatmentTimers[response.UID] = time.AfterFunc(2*time.Hour, func() {
		formattedTime := response.CreatedAt.Format("15:04")
		nextPeriod := response.CreatedAt.Add(2 * time.Hour)
		formattedNextPeriod := nextPeriod.Format("15:04")
		// На болюс 0,4 введенный в 15:00 значение сахара через 2 часа (17:00) 10 ммоль/л
		message := fmt.Sprintf(bolusText, response.Insulin, formattedTime, formattedNextPeriod, sgv)
		c.notifier.Send(message)
	})

	if t, exists := nextTreatmentTimers[response.UID]; exists {
		t.Stop()
	}
	fmt.Println("treatment timer 4 hours is running")
	nextTreatmentTimers[response.UID] = time.AfterFunc(4*time.Hour, func() {
		formattedTime := response.CreatedAt.Format("15:04")
		nextPeriod := response.CreatedAt.Add(4 * time.Hour)
		formattedNextPeriod := nextPeriod.Format("15:04")
		// На болюс 0,4 введенный в 15:00 значение сахара через 4 часа (19:00) 10 ммоль/л
		message := fmt.Sprintf(bolusText, response.Insulin, formattedTime, formattedNextPeriod, sgv)
		c.notifier.Send(message)
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
