package server

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"net/http"
)

type Server struct {
	store *db.UserStore
	Game  *game.Service
	Mux   *http.ServeMux
	WsURL string

	in  chan ClientMessage
	out chan ServerMessage

	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewServer(store *db.UserStore, game *game.Service, url string) *Server {
	return &Server{
		store:      store,
		Game:       game,
		Mux:        http.NewServeMux(),
		WsURL:      url,
		in:         make(chan ClientMessage),
		out:        make(chan ServerMessage),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Start() {
	s.InitRoutes()
	go s.run()
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
