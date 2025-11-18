package app

import (
	"context"
	"errors"
	"feed-api/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(*Application) error

func WithMigrations(dsn, migrationsPath string) Option {
	return func(a *Application) error {
		log.Println("Running migrations...")
		err := repository.RunMigrations(dsn, migrationsPath)
		if err != nil {
			return err
		}
		log.Println("Migrations completed successfully")
		return nil
	}
}

type Application struct {
	server          *http.Server
	worker          Worker
	subscriber      Subscriber
	ctx             context.Context
	cancel          context.CancelFunc
	conn            *pgxpool.Pool
	messageProducer Producer[*repository.Message]
	eventProducer   Producer[*repository.Message]
}

func NewApplication(
	ctx context.Context,
	cancel context.CancelFunc,
	server *http.Server,
	worker Worker,
	subscriber Subscriber,
	conn *pgxpool.Pool,
	messageProducer, eventProducer Producer[*repository.Message],
	options ...Option,
) (*Application, error) {

	app := &Application{
		server:          server,
		worker:          worker,
		subscriber:      subscriber,
		conn:            conn,
		messageProducer: messageProducer,
		eventProducer:   eventProducer,
		ctx:             ctx,
		cancel:          cancel,
	}

	for _, opt := range options {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *Application) Start() error {
	go a.worker.Start(a.ctx)
	go a.subscriber.Run(a.ctx)

	go func() {
		log.Printf("Listening on port %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	return a.awaitShutdown()
}

func (a *Application) awaitShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal, stopping server...")
	a.cancel()
	a.worker.Stop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.server.Shutdown(shutdownCtx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Println("Server shutdown completed successfully")
	}
	a.messageProducer.Close()
	a.eventProducer.Close()
	a.conn.Close()

	log.Println("App stopped gracefully")

	return nil
}
