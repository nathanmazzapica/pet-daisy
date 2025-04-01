package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nathanmazzapica/pet-daisy/logger"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		logger.LogError(fmt.Errorf("failed to connect to database: %w", err))
	}

	fmt.Println("Connected to database!")
}
