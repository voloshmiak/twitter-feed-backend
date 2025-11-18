package main

import (
	"context"
	"feed-bot/internal/bot"
	"feed-bot/internal/client"
	"feed-bot/internal/generator"
	"feed-bot/internal/scheduler"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	UserCount     = 3
	SleepInterval = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	APIEndpoint := fmt.Sprintf("http://%s:%s/api/messages",
		os.Getenv("API_HOST"), os.Getenv("API_PORT"))

	messageGenerator := generator.NewMessageGenerator(UserCount)
	messageFactory := bot.NewMessageFactory()
	httpClient := client.NewHTTPClient[*bot.Message](APIEndpoint)

	messageBot := bot.NewBot[*bot.Message](messageGenerator, messageFactory, httpClient)

	taskScheduler := scheduler.NewScheduler(ctx, cancel, SleepInterval, messageBot)

	err := taskScheduler.Run()
	if err != nil {
		log.Printf("task scheduler error: %s", err)
	}
}
