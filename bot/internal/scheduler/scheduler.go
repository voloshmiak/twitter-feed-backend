package scheduler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Task interface {
	Execute(ctx context.Context) error
}

type Scheduler struct {
	interval time.Duration
	task     Task
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewScheduler(
	ctx context.Context,
	cancel context.CancelFunc,
	interval time.Duration,
	task Task,
) *Scheduler {
	return &Scheduler{
		interval: interval,
		task:     task,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *Scheduler) Run() error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.task.Execute(s.ctx); err != nil {
					continue
				}
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return s.awaitShutdown()
}

func (s *Scheduler) awaitShutdown() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received, stopping scheduler...")
	s.cancel()
	log.Println("Scheduler gracefully stopped")

	return nil
}
