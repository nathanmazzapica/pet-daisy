package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/game"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var WsUrl string

func ServeHome(w http.ResponseWriter, r *http.Request) {

	user_id, err := r.Cookie("user_id_daisy")
	var userID string
	var user *db.User

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):

			user = db.CreateNewUser()
			fmt.Println("hello,", user.DisplayName)
			fmt.Println("newID:", user.UserID)

			domain := ""

			if strings.Contains(r.Host, "pethenry.com") {
				domain = ".pethenry.com"
			}

			cookie := http.Cookie{
				Name:     "user_id_daisy",
				Value:    user.UserID,
				HttpOnly: true,
				Expires:  time.Now().AddDate(10, 0, 0),
				Domain:   domain,
			}
			http.SetCookie(w, &cookie)
		default:
			fmt.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			// todo: make funny error html page
			return
		}
	} else {
		userID = user_id.Value
		user, err = db.GetUserFromDB(userID)
		if err != nil {
			fmt.Println(err)
		}

	}

	fmt.Printf("USER: {%s} CONNECTED\n", user.DisplayName)

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
		TotalPets: game.Counter,
		WS_URL:    WsUrl,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println("error sending html", err)
	}
}

func PostSyncCode(w http.ResponseWriter, r *http.Request) {
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

	userID, err := db.FindIDBySyncCode(data.Code)

	if err != nil {
		fmt.Println("Error recovering user:", err)
		return
	}

	domain := ""

	if strings.Contains(r.Host, "pethenry.com") {
		domain = ".pethenry.com"
	}

	cookie := &http.Cookie{
		Name:     "user_id_daisy",
		Value:    userID,
		HttpOnly: true,
		Expires:  time.Now().AddDate(10, 0, 0),
		Domain:   domain,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"refresh": true})
}

func ServeRoadmap(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/roadmap.html"))
	err := tmpl.Execute(w, nil)

	if err != nil {
		fmt.Println("error sending html", err)
	}
}

func ServeError(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/error.html"))

	_ = tmpl.Execute(w, nil)
}
