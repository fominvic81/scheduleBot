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

	return db, nil
}
