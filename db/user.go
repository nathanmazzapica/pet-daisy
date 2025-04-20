package db

import (
	"sync"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type User struct {
	UserID      string `field:"user_id"`
	DisplayName string `field:"display_name"`
	SyncCode    string `field:"sync_code"`
	PetCount    int    `field:"pets"`
	mu          sync.Mutex
}

func (u *User) ID() string {
	return u.UserID
}

func (u *User) SafeIncrementPet() {
	u.mu.Lock()
	u.PetCount++
	u.mu.Unlock()
}
