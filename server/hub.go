package server

import (
	"encoding/json"
	"github.com/nathanmazzapica/pet-daisy/game"
	"log"
	"strings"
)

type Hub struct {
	clients map[*Client]bool

	receive chan ClientMessage

	broadcast chan ServerMessage

	register chan *Client

	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan ServerMessage),
		receive:    make(chan ClientMessage),
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
			log.Println("Received message:", message)

			// This needs to be moved eventually.
			if strings.Contains(message.Message, "$!pet;") {
				game.PetDaisy(&message.Client.user)
				petCountUpdate := newPetNotification()

				h.broadcast <- petCountUpdate.toBytes()

				// ditto... this is MESSY imo but it works for now
				lbData := game.GetTopX(10)

				data, err := json.Marshal(lbData)
				if err != nil {
					log.Println(err)
					continue
				}

				leaderboardUpdate := ServerMessage{"leaderboard", string(data)}

				h.broadcast <- leaderboardUpdate

				continue
			}

			h.broadcast <- message.toServerMessage()
		}
	}
}

func (h *Hub) broadcastMessages() {
	for {
		message := <-h.broadcast

		log.Printf("Broadcasting message: %s to %d clients", message, len(h.clients))

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
