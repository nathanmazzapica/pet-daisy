package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"

	"github.com/gorilla/websocket"
	_ "net/http/pprof"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients               = make(map[*Client]bool)
	webhookCooldowns      = make(map[string]time.Time)
	mu                    sync.RWMutex
	messages              = make(chan ClientMessage)
	notifications         = make(chan ClientMessage)
	topPlayerCount        = 10
	lastLeaderboardUpdate = int64(0)
)

// Client represents a WebSocket connection
type Client struct {
	conn        *websocket.Conn
	id          string
	user        db.User
	lastPetTime time.Time
	petTimes    [50]time.Time
	sessionPets int
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

	if lastConnectTime, ok := webhookCooldowns[client.user.UserID]; !ok || lastConnectTime.Before(time.Now().Add(-2*time.Minute)) {
		sendDiscordWebhook("🌼 " + client.user.DisplayName + " has connected to Daisy! 🌼")
		webhookCooldowns[client.user.UserID] = time.Now()
	}

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
			clientMsg.Name = client.user.DisplayName
			messages <- clientMsg

			if strings.Contains(strings.ToLower(clientMsg.Message), "daisy") {

				switch strings.ToLower(clientMsg.Message) {
				case "hi daisy":
					messages <- ClientMessage{"Daisy", fmt.Sprintf("hi %v", client.user.DisplayName)}
				case "i love you daisy":
					messages <- ClientMessage{"Daisy", fmt.Sprintf("i love you %v", client.user.DisplayName)}
				case "daisy why are you so cute":
					messages <- ClientMessage{"Daisy", fmt.Sprintf("stop flirting with me %v... arf", client.user.DisplayName)}
				default:
					messages <- daisyMessage()
				}
			}
		}
	}
}

// handlePet increments the pet count safely
func handlePet(client *Client) {
	petTimeIdx := client.sessionPets % 50

	client.petTimes[petTimeIdx] = time.Now()
	if client.sessionPets > 0 && client.sessionPets%10 == 0 {
		//fmt.Println("calculate st dev")
		intervals := make([]int64, 0)

		// I start at 2 because for some reason the first interval is always a large negative. I'll figure it out later
		for i := 2; i < 50; i++ {
			if client.petTimes[i].IsZero() || client.petTimes[i-1].IsZero() {
				continue
			}
			intervals = append(intervals, client.petTimes[i].Sub(client.petTimes[i-1]).Milliseconds())
		}

		//fmt.Println("Intervals:", intervals)
		mean := meanTime(intervals)
		deviation := stdDev(intervals, mean)
		//fmt.Println("Std Dev:", deviation)

		if deviation < 0.5 {
			cheaterCallout := fmt.Sprintf("😡 %s is cheating!! 😡", client.user.DisplayName)
			notifications <- serverNotification(cheaterCallout)
			sendDiscordWebhook(cheaterCallout)

			client.user.PetCount -= 50
			game.Counter -= 50

			client.conn.Close()
			return
		}

	}

	client.sessionPets++

	game.PetDaisy(&client.user)
	client.lastPetTime = time.Now()
	client.user.SaveToDB()

	notifications <- newPetNotification()

	if shouldUpdateLeaderboard() {
		newData := game.GetTopX(topPlayerCount)

		data, err := json.Marshal(newData)
		if err != nil {
			fmt.Println("error encoding json:", err)
		}
		notifications <- ClientMessage{"leaderboard", string(data)}
	}

	personal := client.user.PetCount

	if personal == 10 || personal == 50 || personal == 100 || personal == 500 || personal%1000 == 0 {
		notifications <- newAchievmentNotification(client.user.DisplayName, personal)
		sendDiscordWebhook("🎉 " + client.user.DisplayName + " has pet daisy " + strconv.Itoa(personal) + " times! 🎉")
	}

	if game.Counter%1000 == 0 {
		notifications <- newMilestoneNotification()
	}

	if game.Counter%10000 == 0 {
		sendDiscordWebhook("🎉 " + " Daisy has been pet " + strconv.Itoa(game.Counter) + " times! 🎉")
	}

	if game.Counter == 1_000_000 {
		activeEvent = "milEvent.js"
		notifications <- ClientMessage{Name: "milEvent", Message: activeEvent}
		trophy := fmt.Sprintf("🌼 %s 🌼", client.user.DisplayName)
		sendDiscordWebhook(fmt.Sprintf("🏆 %s was Daisy's 1,000,000th pet! 🏆", client.user.DisplayName))

		client.user.UpdateDisplayName(trophy)
		client.user.DisplayName = trophy

		nameUpdate := ClientMessage{Name: "updateDisplay", Message: trophy}
		sendJSONToClient(client, nameUpdate)

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
		messages <- daisyMessage()
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

func sendDiscordWebhook(message string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	jsonData := []byte(`{"content": "` + message + `"}`)

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending webhook:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Discord webhook returned:", resp.Status)
	}
}

func shouldUpdateLeaderboard() bool {
	now := time.Now().UnixMilli()
	if now > lastLeaderboardUpdate+250 {
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

	err = client.conn.WriteMessage(websocket.TextMessage, jsonData)

	if err != nil {
		fmt.Println("error networking message", err)
		client.conn.Close()
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
	}
}
