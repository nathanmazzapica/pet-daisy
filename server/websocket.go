package server

import (
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"log"
	"strconv"
)

func (s *Server) listen() {
	go s.broadcast()
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

		s.Game.PetDaisy(message.Client.user)

		diff := s.LB.UpdateUser(message.Client.user)
		if len(diff) > 0 {
			s.out <- leaderboardDeltaNotification(diff)
		}

		s.out <- ServerMessage{
			Name: "petCounter",
			Data: strconv.Itoa(int(s.Game.PetCount)),
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

	client.send <- leaderboardUpdateNotification(s.LB.GetAll())

	s.out <- playerJoinNotification(client.DisplayName())
	s.out <- playerCountNotification(len(s.clients))

	utils.SendPlayerConnectionWebhook(client.DisplayName())

}

func (s *Server) handleClientUnregister(client *Client) {
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
		close(client.send)
		s.out <- playerLeftNotification(client.user.DisplayName)
		s.out <- playerCountNotification(len(s.clients))
		s.store.SaveUserScore(client.user)
	}
}

func (s *Server) broadcast() {
	for {
		message := <-s.out

		for client := range s.clients {
			select {
			case client.send <- message:
			default:
				s.unregister <- client
			}
		}
	}
}
