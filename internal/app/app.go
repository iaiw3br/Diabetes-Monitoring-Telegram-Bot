package app

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"notifier/internal/adapters"
	"notifier/internal/usecases"
	"os"
	"time"
)

const (
	highThreshold    = 10.0
	lowThreshold     = 4.5
	noDateDifference = 10 * time.Minute
)

var (
	prevValue           float64
	timer               *time.Timer
	treatments          map[uuid.UUID]bool
	treatmentTimers     map[uuid.UUID]*time.Timer
	nextTreatmentTimers map[uuid.UUID]*time.Timer
	highText            = "Внимание! Сахар выше 10! Текущее значение %.1f ммоль/л"
	lowText             = "Внимание! Сахар ниже 4.5! Текущее значение %.1f ммоль/л"
	bolusText           = "На болюс %.1f введенный в %s значение сахара через 2 часа (%s) %.1f ммоль/л"
	noDataText          = "Нет новых данных"
	interval            = 5 * time.Minute
	//tickerInterval      = 1 * time.Minute
	tickerInterval = 10 * time.Second
)

func Run() {
	fetcher := adapters.NewHttpFetcher()
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	notifier := adapters.NewNotifier(chatID, botToken)

	checker := usecases.NewChecker(fetcher, notifier, interval)

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := checker.CheckAndNotify()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	}

	//response, err := adapters.FetchData()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//treatment, err := adapters.FetchTreatments()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//sgv := response.SGV
	//value := calculateValue(sgv)
	//fmt.Println(value)

	//if isValidDate(treatment.CreatedAt) && treatment.Insulin != 0 {
	//	treatments[treatment.UID] = true
	//
	//	if t, exists := treatmentTimers[treatment.UID]; exists {
	//		t.Stop()
	//	}
	//	fmt.Println("treatment timer 2 hours is running")
	//	treatmentTimers[treatment.UID] = time.AfterFunc(2*time.Hour, func() {
	//		parsedTime, err := parseDate(treatment.CreatedAt)
	//		if err != nil {
	//			fmt.Println("Error parsing time:", err)
	//			return
	//		}
	//		formattedTime := parsedTime.Format("15:04")
	//		nextPeriod := parsedTime.Add(2 * time.Hour)
	//		formattedNextPeriod := nextPeriod.Format("15:04")
	//		// На болюс 0,4 введенный в 15:00 значение сахара через 2 часа (17:00) 10 ммоль/л
	//		message := fmt.Sprintf(bolusText, treatment.Insulin, formattedTime, formattedNextPeriod, value)
	//		sendMessageImmediately(message)
	//	})
	//
	//	if t, exists := nextTreatmentTimers[treatment.UID]; exists {
	//		t.Stop()
	//	}
	//	fmt.Println("treatment timer 4 hours is running")
	//	nextTreatmentTimers[treatment.UID] = time.AfterFunc(4*time.Hour, func() {
	//		parsedTime, err := parseDate(treatment.CreatedAt)
	//		if err != nil {
	//			fmt.Println("Error parsing time:", err)
	//			return
	//		}
	//		formattedTime := parsedTime.Format("15:04")
	//		nextPeriod := parsedTime.Add(4 * time.Hour)
	//		formattedNextPeriod := nextPeriod.Format("15:04")
	//		// На болюс 0,4 введенный в 15:00 значение сахара через 4 часа (19:00) 10 ммоль/л
	//		message := fmt.Sprintf(bolusText, treatment.Insulin, formattedTime, formattedNextPeriod, value)
	//		sendMessageImmediately(message)
	//	})
	//}

	//if isValidDate(response.DateString) {
	//	sendMessageImmediately(noDataText)
	//	return
	//}
	//
	//message := getMessage(value)
	//if message != "" {
	//	sendMessageWithDelay(message, value)
	//}
	//prevValue = value
}

func isValidDate(date time.Time) bool {
	localTime := time.Now().UTC()

	if localTime.Sub(date) > noDateDifference {
		return true
	}
	return false
}

func isValidTimeStamp(date int) bool {
	seconds := date / 1000
	nanoseconds := (date % 1000) * int(time.Millisecond)

	// Create time.Time object in UTC
	timeUTC := time.
		Unix(int64(seconds), int64(nanoseconds)).
		UTC()

	localTime := time.Now().UTC()

	// Print the result
	if localTime.Sub(timeUTC) > 3*time.Minute {
		return true
	}
	return false
}

func parseDate(date string) (*time.Time, error) {
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}
	return &dateTime, nil
}

func getMessage(value float64) string {
	switch {
	case value > highThreshold:
		return fmt.Sprintf(highText, value)
	case value < lowThreshold:
		return fmt.Sprintf(lowText, value)
	case (prevValue > lowThreshold && prevValue < highThreshold) && (value < lowThreshold || value > highThreshold):
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

func sendMessageWithDelay(message string, value float64) {
	if timer != nil {
		timer.Stop()
		timer = nil
	}

	if (prevValue > lowThreshold && prevValue < highThreshold) && !(value > lowThreshold && value < highThreshold) {
		sendMessageImmediately(message)
	} else if !(prevValue > lowThreshold && prevValue < highThreshold) && !(value > lowThreshold && value < highThreshold) {
		timer = time.AfterFunc(5*time.Minute, func() {
			sendMessageImmediately(message)
		})
	}
}

func sendMessageImmediately(message string) {
	fmt.Println("sending message to telegram:", message)
	//botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	//chatID := os.Getenv("TELEGRAM_CHAT_ID")
	//err := adapters.SendTelegramMessage(
	//	botToken,
	//	chatID,
	//	message,
	//)
	//if err != nil {
	//	log.Fatalf("Failed to send message: %v", err)
	//}
}
