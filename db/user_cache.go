package db

import "time"

type UserCache struct {
	Users     map[string]*User
	TempUsers map[string]*User
}

type UserCacheRow struct {
	user   *User
	expiry time.Time
}

func NewUserCache() *UserCache {
	return &UserCache{
		Users:     make(map[string]*User),
		TempUsers: make(map[string]*User),
	}
}

func (c *UserCache) GetUser(userID string) *User {
	if user, ok := c.Users[userID]; ok {
		return user
	} else if user, ok := c.TempUsers[userID]; ok {
		return user
	}
	return nil
}

func (c *UserCache) AddUser(user *User) {
	c.Users[user.UserID] = user
}

func (c *UserCache) AddTempUser(user *User) {
	c.TempUsers[user.UserID] = user
}
