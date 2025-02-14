package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients        = make(map[*Client]bool)
	mu             sync.RWMutex
	messages       = make(chan ClientMessage)
	notifications  = make(chan ClientMessage)
	topPlayerCount = 10
)

// Client represents a WebSocket connection
type Client struct {
	conn        *websocket.Conn
	id          string
	user        db.User
	lastPetTime time.Time
	susPets     int
}

type ClientMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// HandleConnections upgrades HTTP to WebSocket and manages clients
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello!")
	fmt.Println("Cookie header:", r.Header.Get("Cookie"))

	userID, err := db.GetUserID(r)
	if err != nil {
		fmt.Println("Could not retrieve user ID:", err)
		return
	}

	user, err := db.GetUserFromDB(userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{conn: conn, id: userID, user: *user}

	defer func() {
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
		conn.Close()
		notifications <- ClientMessage{"playerCount", strconv.Itoa(len(clients))}
	}()

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	fmt.Println("Client connected.")

	notifications <- newPetNotification()
	notifications <- playerJoinNotification(client.user.DisplayName)
	notifications <- ClientMessage{"playerCount", strconv.Itoa(len(clients))}

	data, err := json.Marshal(game.GetTopX(topPlayerCount))
	notifications <- ClientMessage{"leaderboard", string(data)}

	readMessages(client)
}

// readMessages handles incoming WebSocket messages
func readMessages(client *Client) {
	conn := client.conn
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			client.user.SaveToDB()
			fmt.Println("Client disconnected.")
			notifications <- playerLeftNotification(client.user.DisplayName)
			break
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			fmt.Println("error decoding json:", err)
			continue
		}

		if strings.Contains(clientMsg.Message, "$!pet;") {
			handlePet(client)
		} else {
			messages <- clientMsg
		}
	}
}

// handlePet increments the pet count safely
func handlePet(client *Client) {
	timeSinceLastPet := time.Since(client.lastPetTime)

	if timeSinceLastPet < (15 * time.Millisecond) {
		client.susPets++
		if client.susPets > 15 {
			notifications <- serverNotification(fmt.Sprintf("%s is cheating!!", client.user.DisplayName))
			return
		}
		return
	} else {
		client.susPets = 0
	}

	game.PetDaisy(&client.user)
	client.lastPetTime = time.Now()
	client.user.SaveToDB()

	notifications <- newPetNotification()
	newData := game.GetTopX(topPlayerCount)

	data, err := json.Marshal(newData)
	if err != nil {
		fmt.Println("error encoding json:", err)
	}
	notifications <- ClientMessage{"leaderboard", string(data)}

	personal := client.user.PetCount

	if personal == 10 || personal == 50 || personal == 100 || personal == 500 || personal%1000 == 0 {
		notifications <- newAchievmentNotification(client.user.DisplayName, personal)
	}

	if game.Counter%1000 == 0 {
		notifications <- newMilestoneNotification()
	}
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

func autoSave() {
	for {
		time.Sleep(3 * time.Minute)
		mu.RLock()
		for client := range clients {
			if err := client.user.SaveToDB(); err != nil {
				fmt.Printf("Error saving user %s to db: %v\nWill retry next autosave", client.user.DisplayName, err)
				continue
			}
			fmt.Printf("Saved user %s to db\n", client.user.DisplayName)
		}
		mu.RUnlock()
		fmt.Println("Autosave complete.")
	}
}

func sendJSONToClient(client *Client, notification ClientMessage) {
	jsonData, err := json.Marshal(notification)

	if err != nil {
		fmt.Println("error encoding json:", err)
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, jsonData)

	if err != nil {
		fmt.Println("error networking message", err)
		client.conn.Close()
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
	}
}
