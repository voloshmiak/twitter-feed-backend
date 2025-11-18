package handler

import (
	"context"
	"encoding/json"
	"feed-api/internal/messaging"
	"feed-api/internal/repository"
	"fmt"
	"log"
	"net/http"
)

type FeedHandler struct {
	broadcaster *Broadcaster
	repo        Repository
}

func NewFeedHandler(broadcaster *Broadcaster, repo Repository) *FeedHandler {
	return &FeedHandler{
		broadcaster: broadcaster,
		repo:        repo,
	}
}

func (f *FeedHandler) GetFeed(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	clientChan := make(chan *messaging.Event[*repository.Message], 10)
	f.broadcaster.Register(clientChan)
	defer f.broadcaster.Unregister(clientChan)

	historicalMessages, err := f.repo.GetAllMessages(r.Context())
	if err != nil {
		log.Println("Error fetching historical messages:", err)
	}
	for _, msg := range historicalMessages {
		err = f.writeSseEvent(rw, msg)
		if err != nil {
			log.Println("Error writing sse-event:", err)
		}
	}
	flusher.Flush()

	for {
		select {
		case event := <-clientChan:
			err = event.Process(r.Context(), func(ctx context.Context, msg *repository.Message) error {
				return f.writeSseEvent(rw, msg)
			})
			if err != nil {
				log.Println("Error writing sse event:", err)
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (f *FeedHandler) writeSseEvent(rw http.ResponseWriter, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(rw, "data: %s\n\n", dataBytes)
	if err != nil {
		return err
	}

	return nil
}
