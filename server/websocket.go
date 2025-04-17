package server

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"time"

	"github.com/nathanmazzapica/pet-daisy/game"
	_ "net/http/pprof"
)

const (
	PET_WINDOW    = 25
	SUS_THRESHOLD = 15
)

var (
	lastLeaderboardUpdate = int64(0)
)

// handlePet checks for cheating and increments the pet count
func handlePet(client *Client) {
	petTimeIdx := client.sessionPets % PET_WINDOW

	client.petTimes[petTimeIdx] = time.Now()
	if client.sessionPets > 0 && client.sessionPets%PET_WINDOW == 0 {
		intervals := make([]int64, 0)

		// I start at 2 because for some reason the first interval is always a large negative. I'll figure it out later
		for i := 2; i < PET_WINDOW; i++ {
			if client.petTimes[i].IsZero() || client.petTimes[i-1].IsZero() {
				continue
			}

			interval := client.petTimes[i].Sub(client.petTimes[i-1]).Milliseconds()

			intervals = append(intervals, interval)
		}

		mean := meanTime(intervals)
		deviation := stdDev(intervals, mean)

		if deviation < 1 || client.susPets >= SUS_THRESHOLD {
			fmt.Println("not good")
			kickCheater(client, PET_WINDOW)
			return
		} else {
			client.susPets = 0
		}

	}

	client.sessionPets++

	//game.PetDaisy(&client.user)
	client.lastPetTime = time.Now()
}

func kickCheater(client *Client, penalty int) {
	cheaterCallout := fmt.Sprintf("😡 %s is cheating!! 😡", client.user.DisplayName)
	utils.SendDiscordWebhook(cheaterCallout)

	client.user.PetCount -= penalty
	game.Counter -= int64(penalty)

	client.conn.Close()
}

func (s *Server) autoSave() {
	for {
		time.Sleep(3 * time.Minute)
		s.mu.RLock()
		for client := range s.clients {
			if err := s.store.SaveUserScore(&client.user); err != nil {
				errStr := fmt.Sprintf("Failed to save user %s to db: %v\nWill retry next autosave", client.user.DisplayName, err)
				logger.ErrLog.Println(errStr)
				continue
			}
			fmt.Printf("Saved user %s to db\n", client.user.DisplayName)
		}
		s.mu.RUnlock()
		fmt.Println("Autosave complete.")
	}
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
