package handler

import (
	"context"
	"feed-api/internal/repository"
)

type Eventable interface {
	MarshalJSON() ([]byte, error)
}

type Producer[T Eventable] interface {
	Publish(ctx context.Context, data T) error
	Close() error
}

type Repository interface {
	SaveMessage(ctx context.Context, msg *repository.Message) error
	GetAllMessages(ctx context.Context) ([]*repository.Message, error)
}
