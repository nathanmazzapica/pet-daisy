package db

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"testing"
)

func TestUserStore_CreateUserAndRetrieve(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewUserStore(db)

	created, err := store.CreateUser()

	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	id := created.ID()
	syncCode := created.SyncCode

	retrieved, err := store.GetUserByID(id)

	if err != nil {
		t.Fatalf("error getting user by id: %v", err)
	}

	assertUserEquals(t, created, retrieved)

	retrieved, err = store.GetUserBySyncCode(syncCode)

	if err != nil {
		t.Fatalf("error getting user by sync_code: %v", err)
	}

	assertUserEquals(t, created, retrieved)

}

func TestUserStore_UpdateDisplayName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewUserStore(db)

	created, err := store.CreateUser()

	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	newName := "updateMeVro"
	store.UpdateDisplayName(created, newName)

	if created.DisplayName != newName {
		t.Fatalf("Expected created display name to be updated to %s, got %s", newName, created.DisplayName)
	}

	retrieved, err := store.GetUserByID(created.ID())

	if err != nil {
		t.Fatalf("error getting user by id: %v", err)
	}

	if retrieved.DisplayName != newName {
		t.Fatalf("Expected retrieved display name to be updated to %s, got %s", newName, retrieved.DisplayName)
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

func assertUserEquals(t *testing.T, want, got *User) {
	t.Helper()

	if want == nil || got == nil {
		t.Fatalf("One of the users is nil\nWant: %+v\nGot: %+v", want, got)
	}

	if want.UserID != got.UserID ||
		want.DisplayName != got.DisplayName ||
		want.SyncCode != got.SyncCode ||
		want.PetCount != got.PetCount {
		t.Errorf("Users are not equal.\nWant: %+v\nGot: %+v", want, got)
	}
}
