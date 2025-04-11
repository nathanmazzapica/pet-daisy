package db

import (
	"database/sql"
	"log"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return UserStore{DB: db}
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
	_, err := s.DB.Exec(
		"UPDATE users SET pets = ? WHERE user_id = ?",
		user.PetCount, user.UserID,
	)

	return err
}

func (s *UserStore) GetUserCount() int {
	var count int
	res := s.DB.QueryRow("SELECT COUNT(*) FROM users")

	err := res.Scan(&count)

	if err != nil {
		log.Fatal(err)
		return -1
	}

	return count
}
