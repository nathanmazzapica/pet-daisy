package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"sync"
)

type UserStore struct {
	DB    *sql.DB
	Cache *UserCache
	mu    sync.RWMutex

	newUsers map[string]*User

	LastLeaderboardUpdate int64
}

// UserStoreInterface exists for future purposes and is currently redundant. I plan to eventually move to MySQL and will create a different UserStore type for it that implements this interface.
type UserStoreInterface interface {
	CreateUser() (*User, error)
	PersistUser(*User) error
	SaveUserScore(*User) error
	GetUserCount() (int, error)
	GetUserByID(id string) (*User, error)
	GetUserBySyncCode(syncCode string) (*User, error)
	GetTotalPetCount() (int, error)
	UpdateDisplayName(user *User, displayName string) error
}

func NewUserStore(db *sql.DB) UserStore {
	return UserStore{
		DB:       db,
		Cache:    NewUserCache(),
		mu:       sync.RWMutex{},
		newUsers: make(map[string]*User),
	}
}

func (s *UserStore) CreateUser() (*User, error) {
	user := &User{
		UserID:      uuid.New().String(),
		DisplayName: getRandomDisplayName(),
		SyncCode:    generateSyncCode(),
		PetCount:    0,
	}

	s.Cache.AddUser(user)
	if err := s.PersistUser(user); err != nil {
		return user, err
	}

	return user, nil
}

// TODO: Collapse into s.CreateUser()
func (s *UserStore) CreateTempUser() (*User, error) {
	user := &User{
		UserID:      uuid.New().String(),
		DisplayName: getRandomDisplayName(),
		SyncCode:    generateSyncCode(),
		PetCount:    0,
	}

	s.Cache.AddUser(user)
	s.newUsers[user.UserID] = user

	return user, nil
}

func (s *UserStore) PersistUser(user *User) error {
	_, err := s.DB.Exec(
		"INSERT INTO users (user_id, pets, display_name, sync_code) VALUES (?, ?, ?, ?)",
		user.UserID,
		user.PetCount,
		user.DisplayName,
		user.SyncCode,
	)

	return err
}

func (s *UserStore) SaveUserScore(user *User) error {
	res, err := s.DB.Exec(
		"UPDATE users SET pets = ? WHERE user_id = ?",
		user.PetCount, user.UserID,
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStore) GetUserCount() (int, error) {
	var count int
	res := s.DB.QueryRow("SELECT COUNT(*) FROM users")

	err := res.Scan(&count)

	return count, err
}

func (s *UserStore) GetUserByID(userID string) (*User, error) {

	if user := s.Cache.GetUser(userID); user != nil {
		return user, nil
	}

	user := &User{}

	res := s.DB.QueryRow("SELECT user_id, display_name, sync_code, pets FROM users WHERE user_id = ?", userID)

	if err := res.Scan(&user.UserID, &user.DisplayName, &user.SyncCode, &user.PetCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user with id: %v", userID)
		}
		return nil, err
	}

	s.Cache.AddUser(user)

	return user, nil
}

func (s *UserStore) GetUserBySyncCode(syncCode string) (*User, error) {
	user := &User{}
	res := s.DB.QueryRow("SELECT user_id, display_name, sync_code, pets FROM users WHERE sync_code = ?", syncCode)

	if err := res.Scan(&user.UserID, &user.DisplayName, &user.SyncCode, &user.PetCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user with sync code: %v", syncCode)
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetTotalPetCount() (int, error) {
	var count int
	res := s.DB.QueryRow("SELECT SUM(pets) FROM users")

	err := res.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil

}

func (s *UserStore) UpdateDisplayName(user *User, displayName string) error {
	userID := user.UserID
	_, err := s.DB.Exec("UPDATE users SET display_name = ? WHERE user_id = ?", displayName, userID)

	if err != nil {
		logger.LogError(err)
		return err
	}

	user.DisplayName = displayName

	return nil
}
