package server

import (
	"net/http"
	"strings"
	"time"
)

// GetIdFromCookie retrieves the User ID from the client request's cookie
func GetIdFromCookie(r *http.Request) (string, error) {
	userID, err := r.Cookie("user_id_daisy")

	if err != nil {
		return "", err
	}

	return userID.Value, nil
}

func (s *Server) newIDCookie(r *http.Request, userID string) *http.Cookie {

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

	return cookie
}
