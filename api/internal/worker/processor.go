// api/internal/worker/processor.go
package worker

import (
	"context"
	"feed-api/internal/messaging"
	"log"
)

type Processor[T messaging.Eventable] interface {
	Process(ctx context.Context, event *messaging.Event[T]) error
}

type DatabaseProcessor[T Eventable] struct {
	producer Producer[T]
	repo     Repository[T]
}

func NewDatabaseProcessor[T Eventable](repo Repository[T], producer Producer[T]) *DatabaseProcessor[T] {
	return &DatabaseProcessor[T]{
		repo:     repo,
		producer: producer,
	}
}

func (p *DatabaseProcessor[T]) Process(ctx context.Context, event *messaging.Event[T]) error {
	log.Printf("Processor received message. Saving to database...")

	return event.Process(ctx, func(ctx context.Context, msg T) error {
		if err := p.repo.SaveMessage(ctx, msg); err != nil {
			log.Printf("Failed to save message: %v", err)
			return err
		}

		if err := p.producer.Publish(ctx, msg); err != nil {
			log.Printf("Failed to publish message: %v", err)
			return err
		}

		log.Println("Successfully processed message")
		return nil
	})
}
