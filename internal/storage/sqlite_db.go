// Package storage SQLite implementation.
//
// Last updated: 2026-06-09
package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// OpenDB opens a SQLite database and enables WAL mode.
func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_sync=NORMAL")
	if err != nil {
		return nil, err
	}
	// Set connection limits for SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, nil
}
