package db

import (
	"database/sql"
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/logger"
	_ "modernc.org/sqlite"
)

func Connect() *sql.DB {
	var err error
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		logger.LogError(fmt.Errorf("failed to connect to database: %w", err))
	}

	fmt.Println("Connected to database!")

	return db
}
