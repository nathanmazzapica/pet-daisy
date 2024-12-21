package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

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
}

type ClientMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	fmt.Println("Hello, Daisy!")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/ping", ping)

	go handleChatMessages()
	go handleNotifications()
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("something fucked up")
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {

	user_id, err := r.Cookie("user_id")
	var userID string

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			fmt.Println("hello, new user!")
			newID := uuid.New().String()
			fmt.Println("newID:", newID)
			cookie := http.Cookie{
				Name:     "user_id",
				Value:    newID,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			userID = newID
		default:
			fmt.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			// todo: make funny error html page
			return
		}
	} else {
		userID = user_id.Value
	}

	fmt.Printf("USER: {%s} CONNECTED\n", userID)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err = tmpl.Execute(w, nil)

	if err != nil {
		fmt.Println("error sending html", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello!")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{conn: conn}

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
	notifications <- joinNotification
	notifications <- playerCountUpdate

	readMessages(conn)

}

func readMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
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

// ping is a debug endpoint to test if the server is reachable
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pong")
	w.Write([]byte("pong"))
}
