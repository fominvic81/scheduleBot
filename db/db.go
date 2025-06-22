package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func columnExists(db *sql.DB, table string, column string) (bool, error) {
	res, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))

	if err != nil {
		return false, err
	}

	var name string
	var sink interface{}

	for res.Next() {
		err = res.Scan(
			&sink,
			&name,
			&sink,
			&sink,
			&sink,
			&sink,
		)
		if err != nil {
			_ = res.Close()
			return false, err
		}
		if name == column {
			return true, res.Close()
		}
	}

	return false, res.Close()
}

func createColumnIfNotExists(db *sql.DB, table string, column string, options string) error {
	columnExists, err := columnExists(db, table, column)
	if err != nil {
		return err
	}
	if !columnExists {
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, options))
		if err != nil {
			return err
		}
	}
	return nil
}

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	db.SetMaxOpenConns(1)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id              INTEGER PRIMARY KEY,
		messages        INTEGER NOT NULL DEFAULT 0,
		firstname       TEXT    NOT NULL,
		lastname        TEXT    NULL,
		username        TEXT    NULL,
		faculty         TEXT    NULL,
		educationForm   TEXT    NULL,
		course          TEXT    NULL,
		studyGroup      TEXT    NULL,
		isAdmin         BOOLEAN NOT NULL DEFAULT FALSE,
		keyboardVersion INTEGER NOT NULL DEFAULT 1
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_settings (
		user_id         INTEGER PRIMARY KEY,
		show_groups     INTEGER NOT NULL DEFAULT 1,
		show_teacher    BOOLEAN NOT NULL DEFAULT TRUE,
		hidden_subjects TEXT    NOT NULL DEFAULT "[]"
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS metrics (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id          INTEGER NOT NULL,
		chat_id          INTEGER NOT NULL,
		event_type       INTEGER NOT NULL,
		content          TEXT    NOT NULL,
		media_type       TEXT    NOT NULL,
		media_id         TEXT    NOT NULL,
		album_id         TEXT    NOT NULL,
		reply_to         INTEGER NOT NULL,
		flags            INTEGER NOT NULL,
		created_at       INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS metrics_user_id_index on metrics (user_id)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS metrics_chat_id_index on metrics (chat_id)")
	if err != nil {
		return nil, err
	}

	err = createColumnIfNotExists(db, "users", "banned_until", "INTEGER NOT NULL DEFAULT 0")
	if err != nil {
		return nil, err
	}

	err = createColumnIfNotExists(db, "users", "state", "INTEGER NOT NULL DEFAULT 0")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_searches (
		user_id INTEGER NOT NULL,
		type    INTEGER NOT NULL,
		value   TEXT    NOT NULL,

		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS user_searches_user_id_index on user_searches (user_id)")
	if err != nil {
		return nil, err
	}

	return db, nil
}
