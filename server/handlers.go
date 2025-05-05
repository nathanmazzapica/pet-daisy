package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:8080" || origin == "https://pethenry.com" || origin == "https://www.pethenry.com"
	},
}

// ServeWebsocket upgrades HTTP to WebSocket and manages clients
func (s *Server) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	userID, err := GetIdFromCookie(r)
	log.Println("I am silencing the unused var error for:", userID)

	if err != nil {
		logger.ErrLog.Println("Could not retrieve user ID:", err)
		return
	}

	var user *db.User

	// TODO: Retrieve User

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		logger.ErrLog.Println(err)
		return
	}

	client := s.newClient(conn, user)

	client.hub.register <- client

	fmt.Println("Client connected.")

	go client.writePump()
	go client.readPump()
}

func (s *Server) ServeHome(w http.ResponseWriter, r *http.Request) {

	if !isValidAgent(r.UserAgent()) {
		http.Error(w, "Agent not supported", http.StatusNotImplemented)
		return
	}

	log.Println("serving home page")
	userIdCookie, err := r.Cookie("user_id_daisy")
	var userID string
	var user *db.User

	if err != nil {
		switch {
		// I want to make this its own func at some point
		case errors.Is(err, http.ErrNoCookie):

			user, err = s.store.CreateTempUser()

			if err != nil {
				log.Println("Error creating user:", err)
				return
			}

			fmt.Println("hello,", user.DisplayName)
			fmt.Println("newID:", user.UserID)

			http.SetCookie(w, newIDCookie(r, user.UserID))
		default:
			logger.LogError(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	} else {
		userID = userIdCookie.Value
		user, err = s.store.GetUserByID(userID)
		if err != nil {
			// TODO: handle userID cookie being present but without a matching db record
			logger.LogError(err)
		}
	}

	fmt.Printf("USER: {%s} CONNECTED\n", user.DisplayName)
	fmt.Println("serving home")

	data := struct {
		User      string
		SyncCode  string
		UserPets  int
		TotalPets int64
		WS_URL    string
	}{
		User:      user.DisplayName,
		SyncCode:  user.SyncCode,
		UserPets:  user.PetCount,
		TotalPets: s.Game.PetCount,
		WS_URL:    s.WsURL,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println("error sending html", err)
		logger.LogError(fmt.Errorf("failed to send html: %w", err))
	}
}

func (s *Server) PostSyncCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.store.GetUserBySyncCode(data.Code)

	if err != nil {
		fmt.Println("Error recovering user:", err)
		return
	}

	userID := user.UserID

	http.SetCookie(w, newIDCookie(r, userID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"refresh": true})
}

func ServeRoadmap(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/roadmap.html"))
	err := tmpl.Execute(w, nil)

	if err != nil {
		logger.LogError(fmt.Errorf("failed to send html: %w", err))
	}
}

func ServeBreak(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/break.html"))
	err := tmpl.Execute(w, nil)

	if err != nil {
		logger.LogError(fmt.Errorf("failed to send html: %w", err))
	}
}

func ServeError(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/error.html"))

	_ = tmpl.Execute(w, nil)
}

func isValidAgent(agent string) bool {
	blockedAgents := []string{
		"curl", "wget", "postmanruntime", "python-requests",
		"go-http-client", "java", "libwww-perl", "httpclient",
		"axios", "scrapy", "httpie", "powershell",
	}

	agent = strings.ToLower(agent)

	for _, blockedAgent := range blockedAgents {
		if strings.Contains(agent, blockedAgent) {
			return false
		}
	}

	return true
}
