package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"notifier/internal/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println("Loaded config")

	c := cron.New()
	_, err = c.AddFunc("@every 5s", func() {
		app.Run()
	})
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}

	c.Start()

	// Block forever
	select {}
}
