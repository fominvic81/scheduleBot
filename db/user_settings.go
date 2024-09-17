package db

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type UserSettings struct {
	db             *sql.DB
	UserId         int64
	ShowGroups     int64
	ShowTeacher    bool
	HiddenSubjects []string
}

func (settings *UserSettings) scan(scanner interface{ Scan(src ...any) error }) error {
	hiddenSubjects := ""
	err := scanner.Scan(
		&settings.UserId,
		&settings.ShowGroups,
		&settings.ShowTeacher,
		&hiddenSubjects,
	)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(hiddenSubjects), &settings.HiddenSubjects)
	return err
}

func GetOrCreateUserSettings(db *sql.DB, id int64) (*UserSettings, error) {
	row := db.QueryRow("SELECT * FROM user_settings WHERE user_id = ?", id)
	settings := UserSettings{db: db}

	err := settings.scan(row)

	if errors.Is(err, sql.ErrNoRows) {
		row = settings.db.QueryRow("INSERT INTO user_settings (user_id) VALUES (?) RETURNING *", id)
		err = settings.scan(row)
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (settings *UserSettings) Save() error {
	hiddenSubjects, err := json.Marshal(settings.HiddenSubjects)
	if err != nil {
		return err
	}
	row := settings.db.QueryRow(`UPDATE user_settings SET
			show_groups = ?,
			show_teacher = ?,
			hidden_subjects = ?
		WHERE user_id = ? RETURNING *`,
		settings.ShowGroups,
		settings.ShowTeacher,
		hiddenSubjects,
		settings.UserId,
	)
	err = settings.scan(row)
	return err
}
