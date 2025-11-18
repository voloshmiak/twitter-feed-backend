package worker

import "context"

type Repository[T any] interface {
	SaveMessage(ctx context.Context, msg T) error
	GetAllMessages(ctx context.Context) ([]T, error)
}

type Eventable interface {
	MarshalJSON() ([]byte, error)
}

type Producer[T any] interface {
	Publish(ctx context.Context, data T) error
	Close() error
}
