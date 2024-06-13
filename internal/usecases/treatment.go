package usecases

import (
	"fmt"
	"notifier/internal/entities"
	"time"
)

const (
	bolusText         = "На болюс %.1f введенный в %s значение сахара через %d часа (%s) %.1f ммоль/л (%s)"
	firstPeriodValue  = 2
	secondPeriodValue = 4
	hhMMLayout        = "15:04"
	treatmentsDate    = 20 * time.Minute
	firstPeriod       = firstPeriodValue * time.Hour
	secondPeriod      = secondPeriodValue * time.Hour
	utcPeriod         = 3 * time.Hour
)

var (
	treatmentTimers     map[string]*time.Timer
	nextTreatmentTimers map[string]*time.Timer
	treatmentsBySGV     map[string]float64
)

func (c *Checker) CheckTreatments() error {
	responses, err := c.fetcher.FetchTreatments()
	if err != nil {
		return err
	}
	localTime := time.Now().UTC()

	for _, response := range responses {
		if isLongTimeAgo(response.CreatedAt, localTime, treatmentsDate) || response.Insulin == 0 {
			continue
		}

		if _, ok := treatmentsBySGV[response.UID]; ok {
			continue
		}
		treatmentsBySGV[response.UID] = c.currentSGV

		stopTimerIfExists(treatmentTimers, response.UID)
		createAndSetTimer(treatmentTimers, response.UID, firstPeriod, func() {
			sendNotification(c, response, firstPeriodValue, firstPeriod, response.UID)
		})

		stopTimerIfExists(nextTreatmentTimers, response.UID)
		createAndSetTimer(treatmentTimers, response.UID, secondPeriod, func() {
			sendNotification(c, response, secondPeriodValue, secondPeriod, response.UID)
		})
	}

	return nil
}

func stopTimerIfExists(timerMap map[string]*time.Timer, key string) {
	if t, exists := timerMap[key]; exists {
		t.Stop()
	}
}

func createAndSetTimer(timerMap map[string]*time.Timer, key string, duration time.Duration, fn func()) {
	timerMap[key] = time.AfterFunc(duration, fn)
}

func sendNotification(c *Checker, response entities.TreatmentResponse, periodValue int, duration time.Duration, key string) {
	now := time.Now().Format(hhMMLayout)
	nextPeriod := response.CreatedAt.Add(utcPeriod)
	formattedNextPeriod := nextPeriod.Format(hhMMLayout)
	oldSGV := treatmentsBySGV[response.UID]
	difference := formatDifference(oldSGV, c.currentSGV)

	message := fmt.Sprintf(
		bolusText,
		response.Insulin,
		formattedNextPeriod,
		periodValue,
		now,
		c.currentSGV,
		difference,
	)
	c.notifier.Send(message)

	if duration == secondPeriod {
		delete(treatmentsBySGV, key)
	}
}

func formatDifference(oldSGV, currentSGV float64) string {
	difference := fmt.Sprintf("%.1f", oldSGV-currentSGV)
	if oldSGV-currentSGV >= 0 {
		difference = fmt.Sprintf("+%.1f", oldSGV-currentSGV)
	}
	return difference
}
