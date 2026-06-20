// Package sqlite — адаптеры персистентности (SQLite) для доменных портов.
package sqlite

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "modernc.org/sqlite" // чистый Go-драйвер SQLite, без CGO
)

// Open открывает соединение и применяет схему.
//
// SQLite допускает только одного писателя в момент времени. При
// параллельном сборе снапшотов (горутина на пользователя) это давало
// ошибку "database is locked (SQLITE_BUSY)". Чтобы этого не было:
//   - WAL-журнал: чтение не блокируется записью;
//   - busy_timeout: драйвер ждёт освобождения блокировки, а не падает сразу;
//   - SetMaxOpenConns(1): запись сериализуется на уровне пула Go (надёжно
//     исключает гонки на запись; нагрузка мизерная — раз в час).
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn(path))
	if err != nil {
		return nil, err
	}

	// Одно соединение — запись в SQLite строго последовательная.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

// dsn собирает строку подключения с прагмами WAL и busy_timeout.
// Синтаксис _pragma поддерживается драйвером modernc.org/sqlite.
func dsn(path string) string {
	q := url.Values{}
	q.Add("_pragma", "journal_mode(WAL)")
	q.Add("_pragma", "busy_timeout(5000)") // мс: ждать блокировку до 5 секунд
	q.Add("_pragma", "foreign_keys(ON)")
	return fmt.Sprintf("file:%s?%s", path, q.Encode())
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
