package db

import "time"

type UserCache struct {
	Users map[string]UserCacheRow
}

type UserCacheRow struct {
	user   *User
	expiry time.Time
}

func NewUserCache() *UserCache {
	return &UserCache{
		Users: make(map[string]UserCacheRow),
	}
}

func (c *UserCache) GetUser(userID string) *User {
	if row, ok := c.Users[userID]; ok {
		return row.user
	}
	return nil
}

func (c *UserCache) AddUser(user *User) {
	expiry := time.Now().Add(time.Hour)
	if row, ok := c.Users[user.UserID]; ok {
		row.expiry = expiry
	}
	c.Users[user.UserID] = UserCacheRow{user, expiry}
}

func (c *UserCache) Clean() {
	for userID, row := range c.Users {
		if row.expiry.Before(time.Now()) {
			delete(c.Users, userID)
		}
	}
}
