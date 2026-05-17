package main

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const MaxRows = 3000

func InitDB(path string) (*sql.DB, error) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InsertMessage(db *sql.DB, msg string) error {
	_, err := db.Exec("INSERT INTO messages (message) VALUES (?)", msg)
	return err
}

// Prunes the messages table by keeping the 3000 newest messages, and deleting
// the older ones.
func pruneMessages(db *sql.DB) {
	result, err := db.Exec(`
		DELETE FROM messages
		WHERE id NOT IN (
			SELECT id FROM messages ORDER BY id DESC LIMIT ?
		)
	`, MaxRows)
	if err != nil {
		slog.Error("prune failed", "error", err)
		return
	}

	n, _ := result.RowsAffected()
	if n > 0 {
		slog.Info("pruned messages", "count", n)
	}
}

func StartPruneLoop(db *sql.DB) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		pruneMessages(db)
	}
}
