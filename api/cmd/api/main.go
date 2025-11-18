package main

import (
	"context"
	"feed-api/internal/app"
	"feed-api/internal/handler"
	"feed-api/internal/messaging"
	"feed-api/internal/repository"
	"feed-api/internal/worker"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	eventsToProcessTopic = "events-to-process"
	eventsProcessedTopic = "events-processed"
	groupID              = "message-group"
	migrationsPath       = "migrations"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := os.Getenv("PORT")
	brokers := []string{fmt.Sprintf("%s:%s",
		os.Getenv("KAFKA_HOST"),
		os.Getenv("KAFKA_PORT"))}
	migrateDSN := fmt.Sprintf("cockroachdb://root:@%s:%s/defaultdb?sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
	appDSN := fmt.Sprintf("postgresql://root@%s:%s/defaultdb?sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))

	conn, err := repository.NewConnection(ctx, appDSN)
	if err != nil {
		cancel()
		log.Println("Failed to connect to database:", err)
		return
	}

	messageProducer := messaging.NewProducer[*repository.Message](brokers, eventsProcessedTopic)
	eventProducer := messaging.NewProducer[*repository.Message](brokers, eventsToProcessTopic)
	messageRepository := repository.NewRepository(conn)
	databaseProcessor := worker.NewDatabaseProcessor[*repository.Message](messageRepository, messageProducer)
	messageWorker := worker.NewWorker[*repository.Message](brokers, eventsToProcessTopic, groupID, databaseProcessor)

	broadcaster := handler.NewBroadcaster()
	subscriber := handler.NewSubscriber(brokers, eventsProcessedTopic, groupID, broadcaster)
	router := handler.NewRouter(eventProducer, messageRepository, broadcaster)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	application, err := app.NewApplication(
		ctx, cancel,
		server,
		messageWorker,
		subscriber,
		conn,
		messageProducer,
		eventProducer,
		app.WithMigrations(migrateDSN, migrationsPath),
	)
	if err != nil {
		log.Println("Failed to start application:", err)
		cancel()
		return
	}
	if err = application.Start(); err != nil {
		log.Println("Failed to start application:", err)
		cancel()
		return
	}
}
