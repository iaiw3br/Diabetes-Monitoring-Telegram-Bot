package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"notifier/internal/adapters"
	"notifier/internal/usecases"
	"time"
)

var (
	interval       = 5 * time.Minute
	tickerInterval = 1 * time.Minute
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	fetcher := adapters.NewHttpFetcher()
	notifier := adapters.NewNotifier()
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
}
