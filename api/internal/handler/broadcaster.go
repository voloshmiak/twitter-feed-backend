package handler

import (
	"feed-api/internal/messaging"
	"feed-api/internal/repository"
	"log"
	"sync"
)

type Broadcaster struct {
	mu      sync.RWMutex
	clients map[chan *messaging.Event[*repository.Message]]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan *messaging.Event[*repository.Message]]struct{}),
	}
}

func (b *Broadcaster) Register(client chan *messaging.Event[*repository.Message]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = struct{}{}
	log.Println("Client registered. Total clients:", len(b.clients))
}

func (b *Broadcaster) Unregister(client chan *messaging.Event[*repository.Message]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[client]; ok {
		delete(b.clients, client)
		close(client)
		log.Println("Client unregistered. Total clients:", len(b.clients))
	}
}

func (b *Broadcaster) Broadcast(msg *messaging.Event[*repository.Message]) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for client := range b.clients {
		select {
		case client <- msg:
		default:
		}
	}
}
