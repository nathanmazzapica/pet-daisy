package server

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"net/http"
	"sync"
)

type Server struct {
	store *db.UserStore
	Game  *game.Service
	LB    *db.Leaderboard

	Mux   *http.ServeMux
	WsURL string

	in  chan ClientMessage
	out chan ServerMessage

	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewServer(store *db.UserStore, game *game.Service, url string) *Server {
	lb := db.NewLeaderboard(store.GetTopPlayers())
	return &Server{
		store:      store,
		Game:       game,
		LB:         lb,
		Mux:        http.NewServeMux(),
		WsURL:      url,
		in:         make(chan ClientMessage, 1024),
		out:        make(chan ServerMessage, 1024),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		mu:         sync.RWMutex{},
	}
}

func (s *Server) Start() {
	s.InitRoutes()
	go s.listen()
	go s.store.Autosave()
}

func (s *Server) InitRoutes() {

	s.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//http.HandleFunc("/", ServeBreak)
	s.Mux.HandleFunc("/", s.ServeHome)
	s.Mux.HandleFunc("/sync", s.PostSyncCode)
	s.Mux.HandleFunc("/roadmap", ServeRoadmap)
	s.Mux.HandleFunc("/error", ServeError)

	s.Mux.HandleFunc("/ws", s.ServeWebsocket)
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
