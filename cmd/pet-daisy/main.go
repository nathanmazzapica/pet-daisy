package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"github.com/nathanmazzapica/pet-daisy/server"
	"github.com/nathanmazzapica/pet-daisy/utils"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	logger.InitLog()
	defer logger.CloseLog()

	dbConn := db.Connect()
	store := db.NewUserStore(dbConn)
	game.InitCounter(store)

	hub := server.NewHub()
	wsServer := server.NewServer(hub, store, "ws://localhost:8080/ws")

	wsServer.Start()

	utils.SendDiscordWebhook("daisy is waking up")

	environment := os.Getenv("ENVIRONMENT")
	switch environment {
	case "dev":
		wsServer := server.NewServer(hub, store, "ws://localhost:8080/ws")
		wsServer.Start()
		err = http.ListenAndServe(":8080", nil)
		utils.SendDiscordWebhook(err.Error())
		log.Fatal(err)
	case "prod":
		wsServer := server.NewServer(hub, store, "wss://pethenry.com/ws")
		wsServer.Start()
		go server.RedirectHTTP()
		err = server.StartHTTPS()
		utils.SendDiscordWebhook(err.Error())
		log.Fatal(err)
	default:
		fmt.Println("Invalid environment configuration")
		return
	}

	utils.SendDiscordWebhook("Daisy is going to sleep")
	log.Println("[SHUTDOWN] Something caused an unexpected shutdown")

}
