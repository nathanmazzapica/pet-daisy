package server

import (
	"log"
	"net/http"
)

func InitRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.HandleFunc("/", ServeBreak)
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/sync", PostSyncCode)

	http.HandleFunc("/ws", HandleConnections)
	http.HandleFunc("/roadmap", ServeRoadmap)
	http.HandleFunc("/error", ServeError)

	go handleChatMessages()
	go handleNotifications()
	go autoSave()
	//go dbWorker()
}

func RedirectHTTP() {
	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://pethenry.com", http.StatusMovedPermanently)
	})))
}

func StartHTTPS() error {
	return http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/pethenry.com/fullchain.pem", "/etc/letsencrypt/live/pethenry.com/privkey.pem", nil)
}
