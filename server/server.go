package server

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"net/http"
)

type Server struct {
	Hub   *Hub
	Store *db.UserStore
	mux   *http.ServeMux
	WsURL string
}

var hub *Hub

func NewServer(hub *Hub, store *db.UserStore, url string) *Server {
	return &Server{hub, store, http.NewServeMux(), url}
}

func (s *Server) Start() {
	s.InitRoutes()
	go s.Hub.run()
	go s.autoSave()
}

func (s *Server) InitRoutes() {

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//http.HandleFunc("/", ServeBreak)
	s.mux.HandleFunc("/", s.ServeHome)
	s.mux.HandleFunc("/sync", s.PostSyncCode)
	s.mux.HandleFunc("/roadmap", ServeRoadmap)
	s.mux.HandleFunc("/error", ServeError)

	s.mux.HandleFunc("/ws", s.HandleConnections)
}

func (s *Server) StartHTTPS() error {
	return http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", s.mux)
}

func RedirectHTTP() {
	logger.LogFatalError(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://pethenry.com", http.StatusMovedPermanently)
	})))
}

func StartHTTPS() error {
	return http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", nil)
}
