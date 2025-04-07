package server

import (
	"github.com/nathanmazzapica/pet-daisy/logger"
	"net/http"
)

var hub *Hub

func InitRoutes() {

	hub = newHub()
	go hub.run()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.HandleFunc("/", ServeBreak)
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/sync", PostSyncCode)

	http.HandleFunc("/ws", HandleConnections)
	http.HandleFunc("/roadmap", ServeRoadmap)
	http.HandleFunc("/error", ServeError)

	go autoSave()
}

func RedirectHTTP() {
	logger.LogFatalError(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://pethenry.com", http.StatusMovedPermanently)
	})))
}

func StartHTTPS() error {
	return http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", nil)
}
