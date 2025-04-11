package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"math/rand"
	"net/http"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type User struct {
	UserID      string `field:"user_id"`
	DisplayName string `field:"display_name"`
	SyncCode    string `field:"sync_code"`
	PetCount    int    `field:"pets"`
	exists      bool
}

func (u *User) ID() string {
	return u.UserID
}

// CreateNewUser creates a new user and attempts to save them to the database. If this fails the user is still created, and future database saves will try again
func CreateNewUser() *User {
	userID := uuid.New().String()
	displayName := getRandomDisplayName()

	newUser := User{userID, displayName, generateSyncCode(), 0, false}

	err := newUser.SaveToDB()

	if err != nil {
		fmt.Println("error saving new user to DB: ", err)
	}

	return &newUser

}

// GetUserFromDB retrieves a user from the database as a pointer with the provided userID
func GetUserFromDB(userID string) (*User, error) {
	user := &User{}

	result := DB.QueryRow("SELECT user_id, display_name, sync_code, pets FROM users WHERE user_id=?", userID)

	if err := result.Scan(&user.UserID, &user.DisplayName, &user.SyncCode, &user.PetCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %v not found!", userID)
		}
		return nil, err
	}

	// For migration purposes
	if user.SyncCode == "NEEDCODEPLS" {
		fmt.Printf("user with id %v needs a sync code", userID)
		_, err := DB.Exec("UPDATE users SET sync_code = ? WHERE user_id = ?", generateSyncCode(), userID)
		if err != nil {
			return nil, err
		}
	}

	user.exists = true

	return user, nil
}

func FindIDBySyncCode(code string) (string, error) {
	var userID string

	result := DB.QueryRow("SELECT user_id FROM users WHERE sync_code = ?", code)

	if err := result.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with SyncCode %v not found!", code)
		}
		return "", err
	}

	return userID, nil
}

// SaveToDB saves a user's pets to DB if they exist, otherwise inserts user into DB
func (u *User) SaveToDB() error {

	if DB == nil {
		fmt.Println("db connection is nil")
		return errors.New("db connection is nil")
	}

	if u.exists {
		_, err := DB.Exec("UPDATE users SET pets = ? WHERE user_id = ? ", u.PetCount, u.UserID)

		return err
	}

	fmt.Println("Inserting new user into DB")

	_, err := DB.Exec("INSERT INTO users (user_id, pets, display_name) VALUES (?, ?, ?)", u.UserID, u.PetCount, u.DisplayName)
	u.exists = true

	return err
}

func (u *User) UpdateDisplayName(name string) {
	_, err := DB.Exec("UPDATE users SET display_name = ? WHERE user_id = ?", name, u.UserID)

	if err != nil {
		logger.LogError(fmt.Errorf("failed to update user display name: %w", err))
	}
}

// GetUserID retrieves the User ID from the client request's cookie
func GetUserID(r *http.Request) (string, error) {
	userID, err := r.Cookie("user_id_daisy")

	if err != nil {
		return "", err
	}

	return userID.Value, nil
}

// getRandomZeroNumber returns a random number padded with 0s
func getRandomZeroNumber() string {
	n := rand.Intn(1_000)
	return fmt.Sprintf("%04d", n)
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
		"loving", "hateful", "spiteful", "rude", "polite", "dasterdly", "depressed"}

	nouns := []string{"Dog", "Watermelon", "Crusader", "Lancer", "Envisage", "Frog",
		"Beetle", "Cellphone", "Python", "Lizard", "Butterfly", "Dragon",
		"Automobile", "Cow", "Henry", "Levi", "Array", "Buzzer", "Balloon", "Book",
		"Calendar", "Burrito", "Corgi", "Pencil", "Pen", "Marker", "Bookshelf",
		"Sharpener", "Can", "Lightbulb", "Flower", "Daisy", "Eraser", "Battery",
		"Butter", "Cantaloupe", "Fridge", "Computer", "Programmer", "Kitty", "Barbell", "Bottle", "Toad", "Beryllium", "Consumer", "President", "Orange", "Entity"}

	fmt.Printf("%d\n", len(adjectives)*len(nouns)*1_000)

	adjI := rand.Intn(len(adjectives))
	nounI := rand.Intn(len(nouns))

	return adjectives[adjI] + nouns[nounI] + getRandomZeroNumber()

}
