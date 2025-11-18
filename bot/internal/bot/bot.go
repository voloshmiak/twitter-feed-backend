package bot

import (
	"context"
	"log"
)

type Bot[T Sendable] struct {
	generator MessageGenerator
	factory   MessageFactory[T]
	sender    Sender[T]
}

func NewBot[T Sendable](generator MessageGenerator, factory MessageFactory[T], sender Sender[T]) *Bot[T] {
	return &Bot[T]{
		generator: generator,
		factory:   factory,
		sender:    sender,
	}
}

func (b *Bot[T]) Execute(ctx context.Context) error {
	userID, content := b.generator.Next()
	message := b.factory.Create(userID, content)

	if err := b.sender.Send(ctx, message); err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Message sent: user=%s, content=%s", userID, content)
	return nil
}
