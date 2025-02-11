package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var db *sql.DB
var WS_URL string
var topPlayers []LeaderboardRowData

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients        = make(map[*Client]bool)
	mu             sync.RWMutex
	counter        int
	messages       = make(chan ClientMessage)
	notifications  = make(chan ClientMessage)
	topPlayerCount = 10
)

type Client struct {
	conn        *websocket.Conn
	id          string
	user        User
	lastPetTime time.Time
}

type ClientMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {

	// load DB

	var err error
	db, err = sql.Open("sqlite3", "./data.db")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	result := db.QueryRow("SELECT SUM(pets) FROM users")
	result.Scan(&counter)

	topPlayers = GetTopX(topPlayerCount)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/sync", PostSyncCode)
	//http.HandleFunc("/profile", serveProfile)

	http.HandleFunc("/ws", HandleConnections)

	http.HandleFunc("/ping", ping)

	go handleChatMessages()
	go handleNotifications()
	go autoSave()

	fmt.Println("Hello, Daisy!")

	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	environment := os.Getenv("ENVIRONMENT")

	switch environment {
	case "dev":
		WS_URL = "ws://localhost:8080/ws"
		err = http.ListenAndServe(":8080", nil)
	case "prod":
		WS_URL = "wss://pethenry.com/ws"
		err = http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", nil)

		go func() {
			log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			})))
		}()
	default:
		fmt.Println("Environment variables not detected")
		return
	}

	if err != nil {
		fmt.Println("something messed up, shutting er down.")
		fmt.Println(err)
	}
}

// I am aware this needs reorganization. I will do it later

type SyncData struct {
	Code string `json:"code"`
}

func PostSyncCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data SyncData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userID, err := FindIDBySyncCode(data.Code)

	if err != nil {
		fmt.Println("Error recovering user:", err)
		return
	}

	domain := ""

	if strings.Contains(r.Host, "pethenry.com") {
		domain = ".pethenry.com"
	}

	cookie := &http.Cookie{
		Name:     "user_id_daisy",
		Value:    userID,
		HttpOnly: true,
		Expires:  time.Now().AddDate(10, 0, 0),
		Domain:   domain,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"refresh": true})
}

func ServeHome(w http.ResponseWriter, r *http.Request) {

	user_id, err := r.Cookie("user_id_daisy")
	var userID string
	var user *User

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):

			user = CreateNewUser()
			fmt.Println("hello,", user.displayName)
			fmt.Println("newID:", user.userID)

			domain := ""

			if strings.Contains(r.Host, "pethenry.com") {
				domain = ".pethenry.com"
			}

			cookie := http.Cookie{
				Name:     "user_id_daisy",
				Value:    user.userID,
				HttpOnly: true,
				Expires:  time.Now().AddDate(10, 0, 0),
				Domain:   domain,
			}
			http.SetCookie(w, &cookie)
		default:
			fmt.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			// todo: make funny error html page
			return
		}
	} else {
		userID = user_id.Value
		user, err = GetUserFromDB(userID)
		if err != nil {
			fmt.Println(err)
		}

	}

	fmt.Printf("USER: {%s} CONNECTED\n", user.displayName)

	data := struct {
		User      string
		SyncCode  string
		UserPets  int
		TotalPets int
		WS_URL    string
	}{
		User:      user.displayName,
		SyncCode:  user.syncCode,
		UserPets:  user.petCount,
		TotalPets: counter,
		WS_URL:    WS_URL,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println("error sending html", err)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello!")
	fmt.Println("Cookie header:", r.Header.Get("Cookie"))

	userID, err := GetUserID(r)
	if err != nil {
		fmt.Println("Could not retrieve user ID:", err)
	}

	user, err := GetUserFromDB(userID)
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
		playerCountUpdate := ClientMessage{"playerCount", strconv.Itoa(len(clients))}
		notifications <- playerCountUpdate
	}()

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	fmt.Println("Client connected.")

	var joinNotification ClientMessage

	joinNotification.Name = "server"
	joinNotification.Message = "A new player has joined. Say hi!"

	playerCountUpdate := ClientMessage{"playerCount", strconv.Itoa(len(clients))}

	notifications <- newPetNotification()
	notifications <- playerJoinNotification(client.user.displayName)
	notifications <- playerCountUpdate

	data, err := json.Marshal(topPlayers)

	notifications <- ClientMessage{"leaderboard", string(data)}

	readMessages(client)

}

func readMessages(client *Client) {

	conn := client.conn
	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			if saveError := client.user.SaveToDB(); saveError != nil {
				fmt.Println("Error saving user to db:", saveError)
				// we should still try again
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("Client disconnected.")
				notifications <- playerLeftNotification(client.user.displayName)
				break
			}
			fmt.Println("error reading message:", err)

			break
		}

		fmt.Println(string(msg))

		if len(msg) > 512 {
			fmt.Println("Client message is too large, discarding")
			continue
		}

		var clientMsg ClientMessage

		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			fmt.Println("error decoding json:", err)
			continue
		}

		// lazy workaround to prevent rewriting entire ws system for now
		// The player can change the name sent through js to emulate these features
		// It only affects front end, but I hate fun.
		// I'm letting them send notification messages though, because I don't hate fun that much.

		if clientMsg.Name == "playerCount" || clientMsg.Name == "leaderboard" || clientMsg.Name == "petCounter" {
			continue
		}

		fmt.Println(clientMsg.Name, ":", clientMsg.Message)

		if strings.Contains(clientMsg.Message, "$!pet;") {

			mu.Lock()
			counter++
			mu.Unlock()
			fmt.Println(clientMsg.Name, "pet daisy! She has now been pet: ", counter, "times.")
			client.user.petCount++
			fmt.Println("client.user.petCount:", client.user.petCount)
			fmt.Println(client.user.displayName)

			client.user.SaveToDB()

			notifications <- newPetNotification()

			// This is expensive..... we shouldn't queue the DB every click for every user.
			newData := GetTopX(topPlayerCount)

			data, err := json.Marshal(newData)

			if err != nil {
				fmt.Println("error encoding json:", err)
			}

			notifications <- ClientMessage{Name: "leaderboard", Message: string(data)}

			// temporary

			if shouldSend := leaderboardNeedsUpdate(newData); shouldSend {
				fmt.Println("Leaderboard needs updating, now!")
				// send the updated data
			}

			personal := client.user.petCount

			if personal == 10 || personal == 50 || personal == 100 || personal == 500 || personal%1000 == 0 {
				notifications <- newAchievmentNotification(clientMsg.Name, personal)
			}

			if counter%1000 == 0 {
				notifications <- newMilestoneNotification()
			}

		} else {
			messages <- clientMsg
		}

	}
}

func autoSave() {
	for {
		time.Sleep(3 * time.Minute)
		mu.RLock()
		for client := range clients {
			if err := client.user.SaveToDB(); err != nil {
				fmt.Printf("Error saving user %s to db: %v\nWill retry next autosave", client.user.displayName, err)
				continue
			}
			fmt.Printf("Saved user %s to db\n", client.user.displayName)
		}
		mu.RUnlock()
		fmt.Println("Autosave complete.")
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

// HELPERS //

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

func newPetNotification() ClientMessage {
	return ClientMessage{"petCounter", strconv.Itoa(counter)}
}

func newMilestoneNotification() ClientMessage {
	return ClientMessage{"Daisy", fmt.Sprintf("Yay! I have been pet %v times!", counter)}
}

func newAchievmentNotification(user string, count int) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has pet daisy %v times!", user, count)}
}

func playerJoinNotification(user string) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has joined! say hi!", user)}
}

func playerLeftNotification(user string) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has disconnected :(", user)}
}

// leaderboardNeedsUpdate is a helper function that determines whether we should send the result of GetTopX to the client
// This needs to be fleshed out a little bit
// Should be true if...
// 1. A new player enters top players
// 2. A top player's pet count increases
//
// I want to avoid querying the DB for every pet, it sounds expensive.
// I am going to learn more about SQL before I proceed
func leaderboardNeedsUpdate(newData []LeaderboardRowData) bool {
	for i := 0; i < len(newData); i++ {
		fmt.Println("checking")
		fmt.Println("new: ", newData[i], "old: ", topPlayers[i])
		if newData[i] != topPlayers[i] {
			topPlayers = newData
			return true
		}
	}

	return false
}

// ping is a debug endpoint to test if the server is reachable
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pong")

	top := GetTopX(topPlayerCount)

	data, err := json.Marshal(top)

	if err != nil {
		fmt.Println("error encoding json:", err)
		return
	}

	w.Write(data)
}
