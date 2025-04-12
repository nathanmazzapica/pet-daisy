package db

import (
	"net/http"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type User struct {
	UserID      string `field:"user_id"`
	DisplayName string `field:"display_name"`
	SyncCode    string `field:"sync_code"`
	PetCount    int    `field:"pets"`
}

func (u *User) ID() string {
	return u.UserID
}

//////////////////
// Hey buddy, this shouldn't be here!

// GetUserID retrieves the User ID from the client request's cookie
func GetUserID(r *http.Request) (string, error) {
	userID, err := r.Cookie("user_id_daisy")

	if err != nil {
		return "", err
	}

	return userID.Value, nil
}
