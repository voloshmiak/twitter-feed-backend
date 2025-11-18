package bot

import "context"

type MessageGenerator interface {
	Next() (userID string, content string)
}

type MessageFactory[T Sendable] interface {
	Create(userID, content string) T
}

type Sendable interface {
	MarshalJSON() ([]byte, error)
}

type Sender[T Sendable] interface {
	Send(ctx context.Context, payload T) error
}
