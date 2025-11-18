package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	id        string
	userID    string
	content   string
	createdAt time.Time
}

func NewMessage(userID, content string) *Message {
	return &Message{
		id:        uuid.New().String(),
		userID:    userID,
		content:   content,
		createdAt: time.Now(),
	}
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Payload struct {
		ID        string    `json:"id"`
		UserID    string    `json:"user_id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}
	var decodedPayload Payload
	if err := json.Unmarshal(data, &decodedPayload); err != nil {
		return err
	}
	m.id = decodedPayload.ID
	m.userID = decodedPayload.UserID
	m.content = decodedPayload.Content
	m.createdAt = decodedPayload.CreatedAt
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	type Payload struct {
		ID        string    `json:"id"`
		UserID    string    `json:"user_id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}
	encodedPayload := &Payload{
		ID:        m.id,
		UserID:    m.userID,
		Content:   m.content,
		CreatedAt: m.createdAt,
	}
	return json.Marshal(encodedPayload)
}
