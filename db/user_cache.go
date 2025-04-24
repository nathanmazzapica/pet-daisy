package db

type UserCache struct {
	Users     map[string]*User
	TempUsers map[string]*User
}

func NewUserCache() *UserCache {
	return &UserCache{
		Users:     make(map[string]*User),
		TempUsers: make(map[string]*User),
	}
}

func (c *UserCache) GetUser(userID string) *User {
	return c.Users[userID]
}
