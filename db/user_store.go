package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"log"
	"math/rand"
	"sync"
	"time"
)

type UserStore struct {
	DB *sql.DB

	Cache map[string]*User
	mu    sync.RWMutex

	LastLeaderboardUpdate int64
}

// UserStoreInterface exists for future purposes and is currently redundant. I plan to eventually move to MySQL and will create a different UserStore type for it that implements this interface.
type UserStoreInterface interface {
	CreateUser() (*User, error)
	PersistUser(*User) error
	SaveUserScore(*User) error
	GetUserCount() (int, error)
	GetUserById(id string) (*User, error)
	GetUserBySyncCode(syncCode string) (*User, error)
	GetTotalPetCount() (int, error)
	UpdateDisplayName(user *User, displayName string) error
}

func NewUserStore(db *sql.DB) UserStore {
	return UserStore{
		DB:    db,
		Cache: map[string]*User{},
		mu:    sync.RWMutex{},
	}
}

func (s *UserStore) CreateUser() (*User, error) {
	user := &User{
		UserID:      uuid.New().String(),
		DisplayName: getRandomDisplayName(),
		SyncCode:    generateSyncCode(),
		PetCount:    0,
	}

	if err := s.PersistUser(user); err != nil {
		return user, err
	}

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
	_, err := s.DB.Exec(
		"UPDATE users SET pets = ? WHERE user_id = ?",
		user.PetCount, user.UserID,
	)

	return err
}

func (s *UserStore) GetUserCount() (int, error) {
	var count int
	res := s.DB.QueryRow("SELECT COUNT(*) FROM users")

	err := res.Scan(&count)

	return count, err
}

func (s *UserStore) GetUserByID(userID string) (*User, error) {
	user := &User{}
	res := s.DB.QueryRow("SELECT user_id, display_name, sync_code, pets FROM users WHERE user_id = ?", userID)

	if err := res.Scan(&user.UserID, &user.DisplayName, &user.SyncCode, &user.PetCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user with id: %v", userID)
		}
		return nil, err
	}

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

func (s *UserStore) GetTopPlayers() []LeaderboardRowData {
	var topUsers []LeaderboardRowData

	rows, err := s.DB.Query("SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT 10")

	if err != nil {
		log.Println("Error getting top players:", err)
		return []LeaderboardRowData{}
	}

	position := 1
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.UserID, &user.DisplayName, &user.PetCount)

		topUsers = append(topUsers, UserToLeaderboardRowData(*user, position))
		position++
	}

	s.LastLeaderboardUpdate = time.Now().UnixMilli()

	return topUsers
}

func (s *UserStore) CacheUser(user *User) {
	s.mu.Lock()
	s.Cache[user.UserID] = user
	s.mu.Unlock()
}

func (s *UserStore) GetUserFromCache(userID string) (*User, error) {
	if user, ok := s.Cache[userID]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("no user with id: %v", userID)
}

func (s *UserStore) Autosave() {
	for {
		time.Sleep(3 * time.Minute)
		s.mu.RLock()
		for _, user := range s.Cache {
			err := s.SaveUserScore(user)
			if err != nil {
				log.Printf("save user score error: %v", err)
				log.Printf("user info dump: %+v", user)
				continue
			}
		}
		s.mu.RUnlock()
	}
}

// generateSyncCode generates a random 6 digit 'syncCode' used for account recovery/syncing
func generateSyncCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func getRandomDisplayName() string {
	adjectives := []string{"big", "long", "small", "golden", "yellow", "black",
		"red", "short", "cunning", "silly", "radical", "sluggish",
		"speedy", "humorous", "shy", "scared", "brave", "intelligent", "stupid",
		"orange", "medium", "austere", "gaudy", "ugly", "beautiful", "sexy",
		"intellectual", "philosophical", "charged", "empty", "full",
		"serious", "vengeful", "malignant", "generous", "complacent",
		"ambitious", "lazy", "dull", "sharp", "splendid", "sexy", "cute",
		"loving", "hateful", "spiteful", "rude", "polite", "dastardly", "depressed"}

	nouns := []string{"Dog", "Watermelon", "Crusader", "Lancer", "Envisage", "Frog",
		"Beetle", "Cellphone", "Python", "Lizard", "Butterfly", "Dragon",
		"Automobile", "Cow", "Henry", "Levi", "Array", "Buzzer", "Balloon", "Book",
		"Calendar", "Burrito", "Corgi", "Pencil", "Pen", "Marker", "Bookshelf",
		"Sharpener", "Can", "Lightbulb", "Flower", "Daisy", "Eraser", "Battery",
		"Butter", "Cantaloupe", "Fridge", "Computer", "Programmer", "Kitty", "Barbell", "Bottle", "Toad", "Beryllium", "Consumer", "President", "Orange", "Entity"}

	fmt.Printf("%d\n", len(adjectives)*len(nouns)*1_000)

	adjI := rand.Intn(len(adjectives))
	nounI := rand.Intn(len(nouns))
	num := fmt.Sprintf("%04d", rand.Intn(1_000))

	return adjectives[adjI] + nouns[nounI] + num

}
