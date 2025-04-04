package server

import (
	"encoding/json"
	"github.com/nathanmazzapica/pet-daisy/game"
	"log"
	"strings"
	"sync/atomic"
)

type Hub struct {
	clients map[*Client]bool

	receive chan []byte

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		receive:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	go h.broadcastMessages()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				client.user.SaveToDB()
			}
		case message := <-h.receive:
			log.Println("Received message:", string(message))
			// process message

			var clientMessage ClientMessage

			if err := json.Unmarshal(message, &clientMessage); err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			if strings.Contains(clientMessage.Message, "$!pet;") {
				atomic.AddInt64(&game.Counter, 1)
				petCountUpdate := newPetNotification()

				data, err := json.Marshal(petCountUpdate)

				if err != nil {
					log.Println("Failed to marshal pet count:", err)
				}

				h.broadcast <- data
				continue
			}

			h.broadcast <- message
		}
	}
}

func (h *Hub) broadcastMessages() {
	for {
		message := <-h.broadcast

		log.Printf("Broadcasting message: %s to %d clients", string(message), len(h.clients))

		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
