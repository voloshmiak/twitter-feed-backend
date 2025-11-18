package handler

import (
	"log"
	"net/http"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) CheckHealth(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("OK"))
	if err != nil {
		log.Printf("health check failed: %s\n", err)
		return
	}
}
