package handler

import (
	"feed-api/internal/repository"
	"net/http"
)

func NewRouter(producer Producer[*repository.Message], repo Repository, broadcaster *Broadcaster) http.Handler {
	router := http.NewServeMux()

	healthHandler := NewHealthHandler()
	feedHandler := NewFeedHandler(broadcaster, repo)
	messagesHandler := NewMessageHandler(producer)

	router.HandleFunc("GET /api/health", healthHandler.CheckHealth)
	router.HandleFunc("GET /api/feed", feedHandler.GetFeed)
	router.HandleFunc("POST /api/messages", messagesHandler.AddMessage)

	return router
}
