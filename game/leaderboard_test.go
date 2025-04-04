package game

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/nathanmazzapica/pet-daisy/db"

	_ "modernc.org/sqlite"
)

func setupTestDB() {
	db.DB, _ = sql.Open("sqlite", ":memory:") // Use in-memory DB
	_, err := db.DB.Exec("CREATE TABLE users (user_id TEXT PRIMARY KEY, display_name TEXT, pets INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}

func insertTestUsers(users []db.User) {
	for _, user := range users {
		_, err := db.DB.Exec("INSERT INTO users (user_id, display_name, pets) VALUES (?, ?, ?)", user.UserID, user.DisplayName, user.PetCount)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestPopulateLeaderboardFromDB(t *testing.T) {
	setupTestDB()

	users := []db.User{
		{UserID: "1", DisplayName: "Alice", SyncCode: "CODE1", PetCount: 50},
		{UserID: "2", DisplayName: "Joe", SyncCode: "CODE1", PetCount: 100},
		{UserID: "3", DisplayName: "Charlie", SyncCode: "CODE1", PetCount: 30},
	}
	insertTestUsers(users)

	populateLeaderboardFromDB()

	// Check top player
	if topPlayers[0].DisplayName != "2" || topPlayers[0].PetCount != 100 {
		t.Errorf("Expected top player to be Bob with 100 pets, got %v", topPlayers[0])
	}

	// Check position mapping
	if userPets["1"] != 2 { // Alice should be in position 2
		t.Errorf("Expected Alice to be in position 2, got %d", userPets["1"])
	}
}

func TestUpdateLeaderboard(t *testing.T) {
	setupTestDB()
	users := []db.User{
		{UserID: "1", DisplayName: "Alice", SyncCode: "CODE1", PetCount: 50},
		{UserID: "2", DisplayName: "Joe", SyncCode: "CODE1", PetCount: 100},
		{UserID: "3", DisplayName: "Charlie", SyncCode: "CODE1", PetCount: 30},
	}
	insertTestUsers(users)
	populateLeaderboardFromDB()

	// Alice jumps to 110 pets (should become the new top player)
	alice := db.User{UserID: "1", DisplayName: "Alice", SyncCode: "CODE1", PetCount: 110}
	updateLeaderboard(&alice)

	if topPlayers[0].DisplayName != "1" {
		t.Errorf("Expected Alice to be at the top, but got %v", topPlayers[0])
		fmt.Println(topPlayers)
	}
	if userPets["1"] != 1 {
		t.Errorf("Expected Alice's position to be 1, got %d", userPets["1"])
		fmt.Println(topPlayers)
	}
}
