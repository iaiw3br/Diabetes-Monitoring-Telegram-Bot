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
	treatmentsBySGV = make(map[string]float64)

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
	localTime := time.Now().UTC()

	if isLongTimeAgo(response.DateString, localTime, noDateDifference) {
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

func isLongTimeAgo(date, localTime time.Time, difference time.Duration) bool {
	if localTime.Sub(date) >= difference {
		return true
	}
	return false
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
