package server

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"net/http"
)

type Server struct {
	Hub   *Hub
	store *db.UserStore
	Game  *game.Controller
	Mux   *http.ServeMux
	WsURL string
}

var hub *Hub

func NewServer(hub *Hub, store *db.UserStore, controller *game.Controller, url string) *Server {
	return &Server{hub, store, controller, http.NewServeMux(), url}
}

func (s *Server) Start() {
	s.InitRoutes()
	go s.Hub.run()
	go s.autoSave()
}

func (s *Server) InitRoutes() {

	s.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//http.HandleFunc("/", ServeBreak)
	s.Mux.HandleFunc("/", s.ServeHome)
	s.Mux.HandleFunc("/sync", s.PostSyncCode)
	s.Mux.HandleFunc("/roadmap", ServeRoadmap)
	s.Mux.HandleFunc("/error", ServeError)

	s.Mux.HandleFunc("/ws", s.HandleConnections)
}

func (s *Server) StartHTTPS() error {
	return http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", s.Mux)
}

func RedirectHTTP() {
	logger.LogFatalError(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://pethenry.com", http.StatusMovedPermanently)
	})))
}

func StartHTTPS() error {
	return http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", nil)
}
