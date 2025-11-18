package app

import (
	"context"
)

type Worker interface {
	Start(ctx context.Context)
	Stop()
}

type Subscriber interface {
	Run(ctx context.Context)
}

type Producer[T any] interface {
	Publish(ctx context.Context, data T) error
	Close() error
}
