package server

import (
	"encoding/json"
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"log"
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

	// The cooldown between sending webhook notifications about relevant clients to the discord server
	webhookCooldowns = make(map[string]time.Time)

	mu            sync.RWMutex
	messages      = make(chan ClientMessage)
	notifications = make(chan ClientMessage)

	topPlayerCount        = 10
	lastLeaderboardUpdate = int64(0)
)

// PetEvent is unused
type PetEvent struct {
	User  *db.User
	Count int
}

// HandleConnections upgrades HTTP to WebSocket and manages clients
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello!")
	fmt.Println("Cookie header:", r.Header.Get("Cookie"))

	userID, err := db.GetUserID(r)
	if err != nil {
		logger.ErrLog.Println("Could not retrieve user ID:", err)
		return
	}

	user, err := db.GetUserFromDB(userID)
	if err != nil {
		logger.ErrLog.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		logger.ErrLog.Println(err)
		return
	}

	client := &Client{conn: conn, id: userID, user: *user, hub: hub, send: make(chan []byte, 256)}

	client.hub.register <- client

	fmt.Println("Client connected.")

	go client.writePump()
	go client.readPump()
}

// handlePet increments the pet count safely
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

	game.PetDaisy(&client.user)
	log.Println("Pet Daisy:", client.user.DisplayName)
	log.Println("New count:", game.Counter)
	client.lastPetTime = time.Now()

	client.user.SaveToDB()

}

func kickCheater(client *Client, penalty int) {
	cheaterCallout := fmt.Sprintf("😡 %s is cheating!! 😡", client.user.DisplayName)
	utils.SendDiscordWebhook(cheaterCallout)

	client.user.PetCount -= penalty
	game.Counter -= int64(penalty)

	client.conn.Close()
}

// BROADCAST HANDLERS //
func handleNotifications() {
	for {
		newNotification := <-notifications

		for client := range clients {
			sendJSONToClient(client, newNotification)
		}
	}
}

func handleChatMessages() {
	for {
		newChatMessage := <-messages
		for client := range clients {
			sendJSONToClient(client, newChatMessage)
		}
	}
}

/*func dbWorker() {
	for event := range dbQueue {
		err := event.User.SaveToDB()
		if err != nil {
			log.Println("Error saving to DB:", err)
			continue
		}
		log.Println("Saved to DB:", event.User.DisplayName)
	}
}

*/

func autoSave() {
	for {
		time.Sleep(3 * time.Minute)
		mu.RLock()
		for client := range clients {
			if err := client.user.SaveToDB(); err != nil {
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

func sendJSONToClient(client *Client, notification ClientMessage) {
	jsonData, err := json.Marshal(notification)

	if err != nil {
		fmt.Println("error encoding json:", err)
		return
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	err = client.conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		logger.ErrLog.Printf("Failed to network message to client: %s\n %v", client.conn.RemoteAddr().String(), err)
		client.conn.Close()
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
	}
}
