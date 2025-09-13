package sqlstore

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct{
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlstore.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil{
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS post(
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL,
		created_at DATETIME);
	`)

	if err != nil{
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if _, err = stmt.Exec(); err != nil{
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

