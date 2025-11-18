package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Eventable interface {
	MarshalJSON() ([]byte, error)
}

type Event[T Eventable] struct {
	id        string
	eventType string
	timestamp time.Time
	version   string
	data      T
}

func NewEventMessage[T Eventable](data T) *Event[T] {
	return &Event[T]{
		id:        uuid.New().String(),
		eventType: "message",
		timestamp: time.Now(),
		version:   "1.0",
		data:      data,
	}
}

func (e *Event[T]) Process(ctx context.Context, fn func(ctx context.Context, data T) error) error {
	return fn(ctx, e.data)
}

func (e *Event[T]) MarshalData() ([]byte, error) {
	return json.Marshal(e.data)
}

func (e *Event[T]) UnmarshalJSON(data []byte) error {
	type Payload struct {
		ID        string          `json:"id"`
		EventType string          `json:"event_type"`
		Timestamp time.Time       `json:"timestamp"`
		Version   string          `json:"version"`
		Data      json.RawMessage `json:"data"`
	}
	var decodedPayload Payload

	if err := json.Unmarshal(data, &decodedPayload); err != nil {
		return err
	}

	e.id = decodedPayload.ID
	e.eventType = decodedPayload.EventType
	e.timestamp = decodedPayload.Timestamp
	e.version = decodedPayload.Version

	var dataValue T
	if err := json.Unmarshal(decodedPayload.Data, &dataValue); err != nil {
		return err
	}
	e.data = dataValue

	return nil
}

func (e *Event[T]) MarshalJSON() ([]byte, error) {
	dataBytes, err := json.Marshal(e.data)
	if err != nil {
		return nil, err
	}

	type Payload struct {
		ID        string          `json:"id"`
		EventType string          `json:"event_type"`
		Timestamp time.Time       `json:"timestamp"`
		Version   string          `json:"version"`
		Data      json.RawMessage `json:"data"`
	}
	encodedPayload := Payload{
		ID:        e.id,
		EventType: e.eventType,
		Timestamp: e.timestamp,
		Version:   e.version,
		Data:      json.RawMessage(dataBytes),
	}
	return json.Marshal(encodedPayload)
}
