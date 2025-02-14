package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/sync", PostSyncCode)

	http.HandleFunc("/ws", HandleConnections)
	http.HandleFunc("/roadmap", ServerRoadmap)

	http.HandleFunc("/ping", ping)
}
