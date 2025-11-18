package handler

import (
	"encoding/json"
	"feed-api/internal/repository"
	"net/http"
)

type MessageHandler struct {
	producer Producer[*repository.Message]
}

func NewMessageHandler(p Producer[*repository.Message]) *MessageHandler {
	return &MessageHandler{
		producer: p,
	}
}

func (m *MessageHandler) AddMessage(rw http.ResponseWriter, r *http.Request) {
	type AddMessageRequest struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	var request AddMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	message := repository.NewMessage(request.UserID, request.Content)

	if err := m.producer.Publish(r.Context(), message); err != nil {
		http.Error(rw, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
