// Package sqlite — адаптеры персистентности (SQLite) для доменных портов.
package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite" // чистый Go-драйвер SQLite, без CGO
)

// Open открывает соединение и применяет схему.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS user (
		id       TEXT PRIMARY KEY,   -- x-api-user
		api_key  TEXT NOT NULL,      -- x-api-key
		name     TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS stats_history (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id   TEXT NOT NULL,
		hp        REAL NOT NULL,
		mp        REAL NOT NULL,
		exp       REAL NOT NULL,
		gp        REAL NOT NULL,
		lvl       INTEGER NOT NULL,
		timestamp TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id)
	);
	CREATE INDEX IF NOT EXISTS idx_stats_user_ts
		ON stats_history(user_id, timestamp DESC);`
	_, err := db.Exec(schema)
	return err
}
