package bot

import "encoding/json"

type messageFactory struct{}

func NewMessageFactory() MessageFactory[*Message] {
	return &messageFactory{}
}

func (f *messageFactory) Create(userID, content string) *Message {
	return &Message{
		userID:  userID,
		content: content,
	}
}

type Message struct {
	userID  string
	content string
}

func (m *Message) MarshalJSON() ([]byte, error) {
	type Payload struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	encodedPayload := Payload{
		UserID:  m.userID,
		Content: m.content,
	}

	return json.Marshal(encodedPayload)
}
