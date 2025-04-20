package server

// A lot of this could be moved to websocket.go

import (
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"log"
	"strconv"
)

func (s *Server) run() {
	go s.broadcastMessages()
	for {
		select {
		case client := <-s.register:
			s.handleClientRegister(client)
		case client := <-s.unregister:
			s.handleClientUnregister(client)
		case message := <-s.in:
			s.handleIncomingMessage(message)
		}
	}
}

func (s *Server) handleIncomingMessage(message ClientMessage) {
	log.Println("Received message:", message)

	if message.Data == "$!pet" {

		s.Game.PetDaisy(&message.Client.user)

		s.out <- ServerMessage{
			Name: "petCounter",
			Data: strconv.Itoa(int(s.Game.PetCount)),
		}

		if s.shouldUpdateLeaderboard() {
			s.out <- leaderboardUpdateNotification(s.store.GetTopPlayers())
		}

		count := message.Client.user.PetCount
		if game.CheckPersonalMilestone(count) {
			s.out <- newAchievmentNotification(message.Client.DisplayName(), count)
		}

		return
	}

	s.out <- message.toServerMessage()
}

func (s *Server) handleClientRegister(client *Client) {
	s.clients[client] = true

	s.out <- playerJoinNotification(client.DisplayName())
	s.out <- playerCountNotification(len(s.clients))
	s.out <- leaderboardUpdateNotification(s.store.GetTopPlayers())

	utils.SendPlayerConnectionWebhook(client.DisplayName())

}

func (s *Server) handleClientUnregister(client *Client) {
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
		close(client.send)
		s.out <- playerLeftNotification(client.user.DisplayName)
		s.out <- playerCountNotification(len(s.clients))
	}
}

func (s *Server) broadcastMessages() {
	for {
		message := <-s.out

		for client := range s.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(s.clients, client)
			}
		}
	}
}
