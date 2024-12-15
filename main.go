package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pong")
	w.Write([]byte("pong"))

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
}

type ClientMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var (
	clients       = make(map[*Client]bool)
	mu            sync.Mutex
	counter       int
	messages      = make(chan ClientMessage)
	notifications = make(chan ClientMessage)
)

func NewPetNotification() ClientMessage {
	return ClientMessage{"petCounter", strconv.Itoa(counter)}
}

func NewMilestoneNotification() ClientMessage {
	return ClientMessage{"Daisy", fmt.Sprintf("Yay! I have been pet %v times!", counter)}
}

func NewAchievmentNotification(user string, count int) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has pet daisy %v times!", user, count)}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err := tmpl.Execute(w, nil)

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

	notifications <- NewPetNotification()
	notifications <- joinNotification
	notifications <- playerCountUpdate

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

			notifications <- NewPetNotification()

			personal, err := strconv.Atoi(strings.Split(clientMsg.Message, ";")[1])

			if err != nil {
				fmt.Print("Could not parse %v into a pet count", clientMsg.Message)
			}

			if personal == 10 || personal == 50 || personal == 100 || personal == 500 || personal%1000 == 0 {
				notifications <- NewAchievmentNotification(clientMsg.Name, personal)
			}

			if counter%100 == 0 {
				notifications <- NewMilestoneNotification()
			}

		} else {
			messages <- clientMsg
		}

	}

}

/*

To improve scalability, I could create a unique channel and goroutine handler for each new client

*/

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

func handleNotifications() {
	for {
		newNotification := <-notifications

		for client := range clients {
			sendJSONToClient(client, newNotification)
		}

	}
}

func handleBroadcasts() {
	for {
		newChatMessage := <-messages

		for client := range clients {
			sendJSONToClient(client, newChatMessage)
		}
	}
}

func main() {
	fmt.Println("Hello, Daisy!")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/ping", ping)

	go handleBroadcasts()
	go handleNotifications()
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("something fucked up")
	}
}
