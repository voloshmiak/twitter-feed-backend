package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer[T Eventable] struct {
	writer *kafka.Writer
}

func NewProducer[T Eventable](brokers []string, topic string) *KafkaProducer[T] {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
	return &KafkaProducer[T]{
		writer: writer,
	}
}

func (p *KafkaProducer[T]) Publish(ctx context.Context, data T) error {
	event := NewEventMessage(data)

	encodedEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.id),
		Value: encodedEvent,
	}

	if err = p.writer.WriteMessages(ctx, msg); err != nil {
		log.Println("Failed to publish message:", err)
		return err
	}

	log.Println("Successfully published message")
	return nil
}

func (p *KafkaProducer[T]) Close() error {
	log.Println("Closing kafka writer")
	return p.writer.Close()
}
