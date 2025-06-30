package server

import (
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"log"
	"strconv"
	"time"
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

	s.out <- playerJoinNotification(client.DisplayName())
	s.out <- playerCountNotification(len(s.clients))
	//	s.out <- leaderboardUpdateNotification(s.store.GetTopPlayers())

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

// Need to rethink leaderboard networking to find healthy balance between network usage and crisp realtimeness
// TODO: Send diffs instead of full LB every time
func (s *Server) updateLeaderboard() {
	for {
		time.Sleep(100 * time.Millisecond)
		// TODO: Implement state to pause leaderboard transmission during bulk save to prevent db lock error
		// TODO: use redis or build own in-memory leaderboard. Polling from the db directly every time just causes problems and isn't good practice
		// TODO: stop sending 1kb of data per user every 100ms like cmon BRO THIS IS TRASH MAKE IT NOT
		// this leaderboard is like 98% of our problems
		//s.out <- leaderboardUpdateNotification(s.store.GetTopPlayers())
	}
}
