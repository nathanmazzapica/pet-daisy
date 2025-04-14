package server

import (
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"log"
)

type Hub struct {
	clients map[*Client]bool

	receive chan ClientMessage

	broadcast chan ServerMessage

	register chan *Client

	unregister chan *Client
}

func NewHub() *Hub {
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
			h.handleClientRegister(client)
		case client := <-h.unregister:
			h.handleClientUnregister(client)
		case message := <-h.receive:
			h.handleIncomingMessage(message)
		}
	}
}

func (h *Hub) handleIncomingMessage(message ClientMessage) {
	log.Println("Received message:", message)

	if message.Data == "$!pet" {

		// I will need to refactor handlePet to allow for proper separation of concerns. For now this will optimistically add pets even if the user is detected to be cheating.

		handlePet(message.Client)
		h.broadcast <- newPetNotification()

		if shouldUpdateLeaderboard() {
			h.broadcast <- leaderboardUpdateNotification()
		}

		count := message.Client.user.PetCount
		if game.CheckPersonalMilestone(count) {
			h.broadcast <- newAchievmentNotification(message.Client.DisplayName(), count)
		}

		if game.CheckMilestone() {
			h.broadcast <- newMilestoneNotification()
		}

		return
	}

	h.broadcast <- message.toServerMessage()
}

func (h *Hub) handleClientRegister(client *Client) {
	h.clients[client] = true
	h.broadcast <- playerJoinNotification(client.DisplayName())
	h.broadcast <- playerCountNotification()
	utils.SendPlayerConnectionWebhook(client.DisplayName())

}

func (h *Hub) handleClientUnregister(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		client.user.SaveToDB()

		h.broadcast <- playerLeftNotification(client.user.DisplayName)
		h.broadcast <- playerCountNotification()
	}
}

func (h *Hub) broadcastMessages() {
	for {
		message := <-h.broadcast

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
