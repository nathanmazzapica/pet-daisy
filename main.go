package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/server"
	"log"
	"net/http"
	"os"
	"time"
)

var WS_URL string
var topPlayers []game.LeaderboardRowData

type Client struct {
	conn        *websocket.Conn
	id          string
	user        db.User
	lastPetTime time.Time
	susPets     int
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db.Connect()
	game.InitCounter()

	server.InitRoutes()

	environment := os.Getenv("ENVIRONMENT")
	switch environment {
	case "dev":
		server.WsUrl = "ws://localhost:8080/ws"
		log.Fatal(http.ListenAndServe(":8080", nil))
	case "prod":
		server.WsUrl = "wss://pethenry.com/ws"
		go server.RedirectHTTP()
		log.Fatal(server.StartHTTPS())
	default:
		fmt.Println("Invalid environment configuration")
		return
	}
}

/*
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello!")
	fmt.Println("Cookie header:", r.Header.Get("Cookie"))

	userID, err := db.GetUserID(r)
	if err != nil {
		fmt.Println("Could not retrieve user ID:", err)
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

	data, err := json.Marshal(game.GetTopX(topPlayerCount))

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

		// PET HANDLING; EXPORT TO OWN FUNCTION ONE DAY //
		if strings.Contains(clientMsg.Message, "$!pet;") {

			timeSinceLastPet := time.Now().Sub(client.lastPetTime)

			// TODO: Add check for "uniformity" if user clicks exactly once every 100ms, it's probably cheating
			if timeSinceLastPet < (15 * time.Millisecond) {
				client.susPets++
				if client.susPets > 15 {
					fmt.Println("RAAAAAAAAAGH STOP CHEATING")
					fmt.Println(timeSinceLastPet)
					notifications <- serverNotification(fmt.Sprintf("%v is cheating :/", client.user.displayName))
					return
				}
				continue
			} else {
				client.susPets = 0
			}

			atomic.AddInt64(&counter, 1)
			client.user.petCount++
			client.lastPetTime = time.Now()

			client.user.SaveToDB()

			notifications <- newPetNotification()

			// This is expensive..... we shouldn't queue the DB every click for every user.
			newData := game.GetTopX(topPlayerCount)

			data, err := json.Marshal(newData)

			if err != nil {
				fmt.Println("error encoding json:", err)
			}

			notifications <- ClientMessage{Name: "leaderboard", Message: string(data)}

			// temporary

			//if shouldSend := leaderboardNeedsUpdate(newData); shouldSend {
			//	fmt.Println("Leaderboard needs updating, now!")
			//	// send the updated data
			//}

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

// HELPERS //

// leaderboardNeedsUpdate is a helper function that determines whether we should send the result of GetTopX to the client
// This needs to be fleshed out a little bit
// Should be true if...
// 1. A new player enters top players
// 2. A top player's pet count increases
//
// I want to avoid querying the DB for every pet, it sounds expensive.
// I am going to learn more about SQL before I proceed
func leaderboardNeedsUpdate(newData []game.LeaderboardRowData) bool {
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

	top := game.GetTopX(topPlayerCount)

	data, err := json.Marshal(top)

	if err != nil {
		fmt.Println("error encoding json:", err)
		return
	}

	w.Write(data)
}

*/
