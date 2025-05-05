package db

import "time"

type UserCache struct {
	Rows map[string]UserCacheRow
}

type UserCacheRow struct {
	user   *User
	expiry time.Time
}

func NewUserCache() *UserCache {
	return &UserCache{
		Rows: make(map[string]UserCacheRow),
	}
}

func (c *UserCache) GetUser(userID string) *User {
	if row, ok := c.Rows[userID]; ok {
		return row.user
	}
	return nil
}

func (c *UserCache) AddUser(user *User) {
	expiry := time.Now().Add(time.Hour)
	if row, ok := c.Rows[user.UserID]; ok {
		row.expiry = expiry
	}
	c.Rows[user.UserID] = UserCacheRow{user, expiry}
}

func (c *UserCache) Clean() {
	for userID, row := range c.Rows {
		if row.expiry.Before(time.Now()) {
			delete(c.Rows, userID)
		}
	}
}
