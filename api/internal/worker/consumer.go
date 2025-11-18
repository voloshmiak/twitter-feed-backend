package worker

import (
	"context"
	"encoding/json"
	"errors"
	"feed-api/internal/messaging"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Start(ctx context.Context)
}

type KafkaConsumer[T messaging.Eventable] struct {
	reader    *kafka.Reader
	processor Processor[T]
}

func NewConsumer[T messaging.Eventable](brokers []string, topic, groupID string, processor Processor[T]) *KafkaConsumer[T] {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &KafkaConsumer[T]{
		reader:    reader,
		processor: processor,
	}
}

func (c *KafkaConsumer[T]) Start(ctx context.Context) {
	log.Println("Starting KafkaConsumer...")
	defer c.reader.Close()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("KafkaConsumer context canceled, stopping...")
				return
			}

			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && errors.Is(kafkaErr, kafka.NotCoordinatorForGroup) {
				log.Println("KafkaConsumer not coordinator for group, retrying...")
				time.Sleep(3 * time.Second)
				continue
			}

			log.Printf("Error fetching message: %v", err)
			continue
		}

		var event messaging.Event[T]
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("Error unmarshalling message", err)
			if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
				log.Printf("Error committing poison message: %v", commitErr)
			}
			continue
		}

		if err = c.processor.Process(ctx, &event); err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		if err = c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}
