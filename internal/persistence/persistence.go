// Package persistence provides metadata storage using pure-Go SQLite.
package persistence

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB encapsulates the SQLite database connection.
type DB struct {
	conn *sql.DB
}

// Open initializes the SQLite database and executes schema migrations.
func Open(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("persistence: open: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, err
	}
	return db, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS workspaces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		root_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		workspace_id TEXT NOT NULL,
		provider_id TEXT NOT NULL,
		model_id TEXT NOT NULL,
		mode TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		request TEXT NOT NULL,
		status TEXT NOT NULL,
		result TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := d.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("persistence: migrate: %w", err)
	}
	return nil
}
