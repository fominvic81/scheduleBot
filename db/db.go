package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

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
		created_at       INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP
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

	return db, nil
}
