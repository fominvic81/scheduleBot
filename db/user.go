package db

import (
	"database/sql"
)

type User struct {
	db              *sql.DB
	Id              int64
	Messages        int64
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
}

func (user *User) scan(scanner interface{ Scan(src ...any) error }) error {
	return scanner.Scan(
		&user.Id,
		&user.Messages,
		&user.Firstname,
		&user.Lastname,
		&user.Username,
		&user.Faculty,
		&user.EducationForm,
		&user.Course,
		&user.StudyGroup,
		&user.IsAdmin,
		&user.KeyboardVersion,
	)
}

func GetOrCreateUser(db *sql.DB, id int64, firstname string) (*User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	user := User{db: db}

	err := user.scan(row)

	if err == sql.ErrNoRows {
		row = user.db.QueryRow("INSERT INTO users (id, firstname) VALUES (?, ?) RETURNING *",
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
			firstname = ?,
			lastname = ?,
			username = ?,
			faculty = ?,
			educationForm = ?,
			course = ?,
			studyGroup = ?,
			isAdmin = ?,
			keyboardVersion = ?
		WHERE id = ? RETURNING *`,
		user.Messages,
		user.Firstname,
		user.Lastname,
		user.Username,
		user.Faculty,
		user.EducationForm,
		user.Course,
		user.StudyGroup,
		user.IsAdmin,
		user.KeyboardVersion,
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
	rows, err := db.Query("SELECT * FROM users WHERE isAdmin = TRUE")
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
