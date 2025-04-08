package utils

import (
	"bytes"
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"net/http"
	"os"
	"time"
)

var cooldowns = make(map[string]time.Time)

func SendDiscordWebhook(message string) {
	if os.Getenv("ENVIRONMENT") == "dev" {
		fmt.Println("not sending discord webhook in dev mode")
		return
	}
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	jsonData := []byte(`{"content": "` + message + `"}`)

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.ErrLog.Println("Failed to create Discord webhook request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.ErrLog.Println("Failed to send webhook:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Discord webhook returned:", resp.Status)
	}
}

func SendPlayerConnectionWebhook(username string) {
	if cooldowns[username].After(time.Now()) {
		return
	}
	message := fmt.Sprintf("%s has connected! 🌼", username)
	SendDiscordWebhook(message)
	cooldowns[username] = time.Now().Add(time.Minute * 5)
}

func SendPlayerDisconnectWebhook(username string) {
	if cooldowns[username].After(time.Now()) {
		return
	}
	message := fmt.Sprintf("%s has disconnected... 💔", username)
	SendDiscordWebhook(message)
	cooldowns[username] = time.Now().Add(time.Minute * 1)
}
