package server

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	_ "net/http/pprof"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	PET_WINDOW    = 25
	SUS_THRESHOLD = 15
)

var (
	// Set of clients
	clients = make(map[*Client]bool)

	mu                    sync.RWMutex
	lastLeaderboardUpdate = int64(0)
)

// PetEvent is unused
type PetEvent struct {
	User  *db.User
	Count int
}

// HandleConnections upgrades HTTP to WebSocket and manages clients
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	userID, err := db.GetUserID(r)
	if err != nil {
		logger.ErrLog.Println("Could not retrieve user ID:", err)
		return
	}

	user, err := s.Store.GetUserByID(userID)
	if err != nil {
		logger.ErrLog.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		logger.ErrLog.Println(err)
		return
	}

	client := &Client{conn: conn, id: userID, user: *user, hub: s.Hub, send: make(chan ServerMessage, 256)}

	client.hub.register <- client

	client.hub.broadcast <- leaderboardUpdateNotification()

	fmt.Println("Client connected.")

	go client.writePump()
	go client.readPump()
}

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
		mu.RLock()
		for client := range clients {
			if err := s.Store.SaveUserScore(&client.user); err != nil {
				errStr := fmt.Sprintf("Failed to save user %s to db: %v\nWill retry next autosave", client.user.DisplayName, err)
				logger.ErrLog.Println(errStr)
				continue
			}
			fmt.Printf("Saved user %s to db\n", client.user.DisplayName)
		}
		mu.RUnlock()
		fmt.Println("Autosave complete.")
	}
}

func getDelay() int64 {
	delay := int64(150*len(clients)) / 2

	if delay > 1000 {
		return 1000
	}

	return delay
}

func shouldUpdateLeaderboard() bool {
	now := time.Now().UnixMilli()
	if now > lastLeaderboardUpdate+getDelay() {
		lastLeaderboardUpdate = now
		return true
	}
	return false
}
