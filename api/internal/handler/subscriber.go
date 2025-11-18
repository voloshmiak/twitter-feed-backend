package handler

import (
	"context"
	"encoding/json"
	"errors"
	"feed-api/internal/messaging"
	"feed-api/internal/repository"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Subscriber struct {
	reader      *kafka.Reader
	broadcaster *Broadcaster
}

func NewSubscriber(brokers []string, topic, groupID string, b *Broadcaster) *Subscriber {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Subscriber{
		reader:      reader,
		broadcaster: b,
	}
}

func (s *Subscriber) Run(ctx context.Context) {
	log.Println("Notification subscriber started")
	defer s.reader.Close()

	for {
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Notification subscriber stopped")
				return
			}

			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && errors.Is(kafkaErr, kafka.NotCoordinatorForGroup) {
				log.Println("KafkaConsumer not coordinator for group, retrying...")
				time.Sleep(3 * time.Second)
				continue
			}

			log.Println("Error fetching message:", err)
			continue
		}

		var event messaging.Event[*repository.Message]
		if err = json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Error unmarshalling message:", err)
			if commitErr := s.reader.CommitMessages(ctx, msg); commitErr != nil {
				log.Println("Error committing poison message:", commitErr)
			}
			continue
		}

		log.Println("Notification subscriber received message")

		s.broadcaster.Broadcast(&event)

		err = s.reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Println("Error committing message:", err)
		}
	}
}
