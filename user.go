package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

type User struct {
	userID      string `field:"user_id"`
	displayName string `field:"display_name"`
	petCount    int    `field:"pets"`
	exists      bool
}

func (u *User) ID() string {
	return u.userID
}

// CreateNewUser creates a new user and attempts to save them to the database. If this fails the user is still created, and future database saves will try again
func CreateNewUser() *User {
	userID := uuid.New().String()
	displayName := getRandomDisplayName()

	newUser := User{userID, displayName, 0, false}

	err := newUser.SaveToDB()

	if err != nil {
		fmt.Println("error saving new user to DB: ", err)
	}

	return &newUser

}

// GetUserFromDB retrieves a user from the database as a pointer with the provided userID
func GetUserFromDB(userID string) (*User, error) {
	user := &User{}

	result := db.QueryRow("SELECT user_id, display_name, pets FROM users WHERE user_id=?", userID)

	if err := result.Scan(&user.userID, &user.displayName, &user.petCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %v not found!", userID)
		}
		return nil, err
	}

	user.exists = true

	return user, nil
}

// im wondering if a key as param would work better here

// SaveToDB saves a user's pets to DB if they exist, otherwise inserts user into DB
func (u *User) SaveToDB() error {

	if db == nil {
		fmt.Println("db connection is nil")
		return errors.New("db connection is nil")
	}

	if u.exists {
		fmt.Println("user already exists, saving pets")
		_, err := db.Exec("UPDATE users SET pets = ? WHERE user_id = ? ", u.petCount, u.userID)

		return err
	}

	fmt.Println("Inserting new user into DB")

	_, err := db.Exec("INSERT INTO users (user_id, pets, display_name) VALUES (?, ?, ?)", u.userID, u.petCount, u.displayName)
	u.exists = true

	return err
}

// getRandomZeroNumber returns a random number padded with 0s
func getRandomZeroNumber() string {
	n := rand.Intn(1_000)
	return fmt.Sprintf("%04d", n)
}

func getRandomDisplayName() string {
	adjectives := []string{"big", "long", "small", "golden", "yellow", "black",
		"red", "short", "cunning", "silly", "radical", "sluggish",
		"speedy", "humorous", "shy", "scared", "brave", "intelligent", "stupid",
		"orange", "medium", "austere", "gaudy", "ugly", "beautiful", "sexy",
		"intellectual", "philosophical", "charged", "empty", "full",
		"serious", "vengeful", "malignant", "generous", "complacent",
		"ambitious", "lazy", "dull", "sharp", "splendid", "sexy", "cute",
		"loving", "hateful", "spiteful", "rude", "polite", "dasterdly"}

	nouns := []string{"Dog", "Watermelon", "Crusader", "Lancer", "Envisage", "Frog",
		"Beetle", "Cellphone", "Python", "Lizard", "Butterfly", "Dragon",
		"Automobile", "Cow", "Henry", "Levi", "Array", "Buzzer", "Balloon", "Book",
		"Calendar", "Burrito", "Corgi", "Pencil", "Pen", "Marker", "Bookshelf",
		"Sharpener", "Can", "Lightbulb", "Flower", "Daisy", "Eraser", "Battery",
		"Butter", "Cantaloupe", "Fridge", "Computer", "Programmer", "Kitty"}

	fmt.Printf("%d\n", len(adjectives)*len(nouns)*1_000)

	adjI := rand.Intn(len(adjectives))
	nounI := rand.Intn(len(nouns))

	return adjectives[adjI] + nouns[nounI] + getRandomZeroNumber()

}
