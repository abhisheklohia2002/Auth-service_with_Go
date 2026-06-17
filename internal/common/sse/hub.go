package sse

import (
	"sync"
)

type ClientChan chan []byte

type Hub struct {
	clients map[ClientChan]bool
	mutex   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[ClientChan]bool),
	}
}

func (h *Hub) AddClient(client ClientChan) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true
}

func (h *Hub) RemoveClient(client ClientChan) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.clients, client)
	close(client)
}

func (h *Hub) Broadcast(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client <- message:
		default:

		}
	}
}