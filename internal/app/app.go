package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"notifier/internal/adapters"
	"notifier/internal/usecases"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	interval       = 5 * time.Minute
	tickerInterval = 30 * time.Second
	runBotText     = "бот включен"
	stopBotText    = "бот отключен"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	fetcher := adapters.NewHttpFetcher()
	notifier := adapters.NewNotifier()
	checker := usecases.NewChecker(fetcher, notifier, interval)
	notifier.Send(runBotText)

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := checker.CheckAndNotify()
				if err != nil {
					fmt.Println("Error:", err)
				}
			case <-ctx.Done():
				fmt.Println("gracefully shutting down")
				return
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	cancel()

	notifier.Send(stopBotText)
	wg.Wait()
	fmt.Println("Application stopped")
}
