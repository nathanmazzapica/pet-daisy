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

// TODO: Refactor!!!
func main() {
	if err := run(); err != nil {
		utils.SendDiscordWebhook(err.Error())
		log.Fatalf("error running server: %v", err)
	}
}

// run runs the server and returns the exit status
func run() error {
	loadEnv()
	logger.InitLog()
	defer logger.CloseLog()

	store := initStore()
	game := game.NewController(&store)

	wsServer := server.NewServer(&store, game, getWebsocketURL())
	wsServer.Start()

	utils.SendDiscordWebhook("daisy is waking up!")

	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		return http.ListenAndServe(":8080", wsServer.Mux)
	case "prod":
		go server.RedirectHTTP()
		return server.StartHTTPS()
	default:
		return fmt.Errorf("invalid environment configuration")
	}

	return fmt.Errorf("not implemented")
}

func initStore() db.UserStore {
	dbConn := db.Connect()
	return db.NewUserStore(dbConn)
}

func loadEnv() error {
	return godotenv.Load()
}

func getWebsocketURL() string {
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		return "ws://localhost:8080/ws"
	case "prod":
		return "wss://pethenry.com/ws"
	default:
		return "ws://localhost:8080/ws"
	}
}
