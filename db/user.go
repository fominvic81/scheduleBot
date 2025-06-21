package db

import (
	"database/sql"
	"errors"
	"time"
)

const (
	UserStateNone = iota
	UserStateSearchGroup
	UserStateSearchTeacher
)

type User struct {
	db              *sql.DB
	Id              int64
	Messages        int64
	State           int64
	Firstname       string
	Lastname        *string
	Username        *string
	Faculty         *string
	EducationForm   *string
	Course          *string
	StudyGroup      *string
	IsAdmin         bool
	KeyboardVersion int64
	Settings        *UserSettings
	BannedUntil     time.Time
}

func (user *User) scan(scanner interface{ Scan(src ...any) error }) error {
	bannedUntil := int64(0)
	err := scanner.Scan(
		&user.Id,
		&user.Messages,
		&user.State,
		&user.Firstname,
		&user.Lastname,
		&user.Username,
		&user.Faculty,
		&user.EducationForm,
		&user.Course,
		&user.StudyGroup,
		&user.IsAdmin,
		&user.KeyboardVersion,
		&bannedUntil,
	)
	user.BannedUntil = time.Unix(bannedUntil, 0).Local()
	return err
}

func GetOrCreateUser(db *sql.DB, id int64, firstname string) (*User, error) {
	row := db.QueryRow(`SELECT 
		id,
		messages,
		state,
		firstname,
		lastname,
		username,
		faculty,
		educationForm,
		course,
		studyGroup,
		isAdmin,
		keyboardVersion,
		banned_until
	FROM users WHERE id = ?`, id)

	user := User{db: db}

	err := user.scan(row)

	if errors.Is(err, sql.ErrNoRows) {
		row = user.db.QueryRow(`
			INSERT INTO users (id, firstname) VALUES (?, ?) RETURNING
				id,
				messages,
				state,
				firstname,
				lastname,
				username,
				faculty,
				educationForm,
				course,
				studyGroup,
				isAdmin,
				keyboardVersion,
				banned_until
			`,
			id,
			firstname,
		)
		err = user.scan(row)
	}
	if err != nil {
		return nil, err
	}

	settings, err := GetOrCreateUserSettings(db, id)
	if err != nil {
		return nil, err
	}

	user.Settings = settings

	return &user, nil
}

func (user *User) Save() error {
	row := user.db.QueryRow(`UPDATE users SET
			messages = ?,
			state = ?,
			firstname = ?,
			lastname = ?,
			username = ?,
			faculty = ?,
			educationForm = ?,
			course = ?,
			studyGroup = ?,
			isAdmin = ?,
			keyboardVersion = ?,
			banned_until = ?
		WHERE id = ? RETURNING
			id,
			messages,
			state,
			firstname,
			lastname,
			username,
			faculty,
			educationForm,
			course,
			studyGroup,
			isAdmin,
			keyboardVersion,
			banned_until
		`,
		user.Messages,
		user.State,
		user.Firstname,
		user.Lastname,
		user.Username,
		user.Faculty,
		user.EducationForm,
		user.Course,
		user.StudyGroup,
		user.IsAdmin,
		user.KeyboardVersion,
		user.BannedUntil.Unix(),
		user.Id,
	)
	err := user.scan(row)
	if err != nil {
		return err
	}

	err = user.Settings.Save()

	return err
}

func GetAdminUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT 
		id,
		messages,
		state,
		firstname,
		lastname,
		username,
		faculty,
		educationForm,
		course,
		studyGroup,
		isAdmin,
		keyboardVersion,
		banned_until
 	FROM users WHERE isAdmin = TRUE`)

	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for rows.Next() {
		user := User{db: db}
		err := user.scan(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
