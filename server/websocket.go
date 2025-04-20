package server

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"time"

	"github.com/nathanmazzapica/pet-daisy/game"
	_ "net/http/pprof"
)

// Deprecated: Controls for pet validation
const (
	PET_WINDOW    = 25
	SUS_THRESHOLD = 15
)

var (
	lastLeaderboardUpdate = int64(0)
)

func kickCheater(client *Client, penalty int) {
	cheaterCallout := fmt.Sprintf("😡 %s is cheating!! 😡", client.user.DisplayName)
	utils.SendDiscordWebhook(cheaterCallout)

	client.user.PetCount -= penalty
	game.Counter -= int64(penalty)

	client.conn.Close()
}

func (s *Server) getDelay() int64 {
	delay := int64(150*len(s.clients)) / 2

	if delay > 1000 {
		return 1000
	}

	return delay
}

func (s *Server) shouldUpdateLeaderboard() bool {
	now := time.Now().UnixMilli()
	if now > lastLeaderboardUpdate+s.getDelay() {
		lastLeaderboardUpdate = now
		return true
	}
	return false
}
