package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var db *sql.DB

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients       = make(map[*Client]bool)
	mu            sync.Mutex
	counter       int
	messages      = make(chan ClientMessage)
	notifications = make(chan ClientMessage)
)

type Client struct {
	conn *websocket.Conn
	id   string
	user User
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/ping", ping)

	go handleChatMessages()
	go handleNotifications()

	fmt.Println("Hello, Daisy!")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("something messed up, shutting er down.")
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {

	user_id, err := r.Cookie("user_id")
	var userID string
	var user *User

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):

			user = CreateNewUser()
			fmt.Println("hello,", user.displayName)
			fmt.Println("newID:", user.userID)
			cookie := http.Cookie{
				Name:     "user_id",
				Value:    user.userID,
				HttpOnly: true,
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
		UserPets  int
		TotalPets int
	}{
		User:      user.displayName,
		UserPets:  user.petCount,
		TotalPets: counter,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println("error sending html", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello!")

	userID, err := getUserID(r)
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

	readMessages(client)

}

func readMessages(client *Client) {

	conn := client.conn
	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("Client disconnected.")
				notifications <- playerLeftNotification(client.user.displayName)
				break
			}
			fmt.Println("error reading message:", err)
			break
		}

		fmt.Println(string(msg))

		var clientMsg ClientMessage

		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			fmt.Println("error decoding json:", err)
			continue
		}

		fmt.Println(clientMsg.Name, ":", clientMsg.Message)

		if strings.Contains(clientMsg.Message, "$!pet;") {

			mu.Lock()
			counter++
			mu.Unlock()
			fmt.Println(clientMsg.Name, "pet daisy! She has now been pet: ", counter, "times.")
			fmt.Println(client.user.displayName)

			notifications <- newPetNotification()

			personal, err := strconv.Atoi(strings.Split(clientMsg.Message, ";")[1])

			if err != nil {
				fmt.Print("Could not parse %v into a pet count", clientMsg.Message)
			}

			if personal == 10 || personal == 50 || personal == 100 || personal == 500 || personal%1000 == 0 {
				notifications <- newAchievmentNotification(clientMsg.Name, personal)
			}

			if counter%100 == 0 {
				notifications <- newMilestoneNotification()
			}

		} else {
			messages <- clientMsg
		}

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

// getUserID retrieves the User ID from the client request's cookie
func getUserID(r *http.Request) (string, error) {
	userID, err := r.Cookie("user_id")

	if err != nil {
		return "", err
	}

	return userID.Value, nil
}

// ping is a debug endpoint to test if the server is reachable
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pong")
	w.Write([]byte("pong"))
}
