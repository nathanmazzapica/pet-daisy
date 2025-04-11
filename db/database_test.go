package db

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	_ "modernc.org/sqlite"
	"testing"
)

func TestUser_SaveToDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewUserStore(db)

	testUser := createDummyUser()

	err := store.PersistUser(testUser)

	if err != nil {
		t.Fatal(err)
	}

	userCount := store.GetUserCount()

	if userCount != 1 {
		t.Errorf("User count should be 1, was %d", userCount)
		t.FailNow()
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE users (user_id TEXT PRIMARY KEY, display_name TEXT, pets INTEGER, sync_code TEXT NOT NULL DEFAULT 'NEEDCODEPLS')")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createDummyUser() *User {
	userID := uuid.New().String()
	displayName := getRandomDisplayName()
	newUser := User{UserID: userID, DisplayName: displayName, SyncCode: generateSyncCode(), PetCount: 0, exists: false}

	return &newUser
}
