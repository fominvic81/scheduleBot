package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	UserStateNone = iota
	UserStateSearchGroup
	UserStateSearchTeacher
)

const (
	UserSearchTypeEmployee = iota
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

func GetUser(db *sql.DB, id int64) (*User, error) {
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

	if err := user.scan(row); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
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
	FROM users WHERE username = ?`, username)

	user := User{db: db}

	if err := user.scan(row); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetOrCreateUser(db *sql.DB, id int64, firstname string) (*User, error) {
	user, err := GetUser(db, id)

	if errors.Is(err, sql.ErrNoRows) {
		row := db.QueryRow(`
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
		user = &User{db: db}
		if err = user.scan(row); err != nil {
			return nil, err
		}
	}

	settings, err := GetOrCreateUserSettings(db, id)
	if err != nil {
		return nil, err
	}

	user.Settings = settings

	return user, nil
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

func (user *User) GetSearches(searchType int) ([]string, error) {
	rows, err := user.db.Query(`SELECT value FROM user_searches WHERE user_id = ? AND type = ?`, user.Id, searchType)
	if err != nil {
		return []string{}, err
	}
	values := []string{}
	for rows.Next() {
		value := ""
		if err := rows.Scan(&value); err != nil {
			return []string{}, err
		}
		values = append(values, value)
	}
	return values, nil
}

func (user *User) SetSearches(searchType int, values []string) error {
	tx, err := user.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM user_searches WHERE user_id = ? AND type = ?", user.Id, searchType); err != nil {
		return err
	}

	if len(values) > 0 {
		insertQuery := "INSERT INTO user_searches (user_id, type, value) VALUES "
		insertValues := []any{}
		for _, value := range values {
			insertQuery += "(?,?,?),"
			insertValues = append(insertValues, fmt.Sprint(user.Id), fmt.Sprint(searchType), value)
		}
		insertQuery = insertQuery[0 : len(insertQuery)-1]

		if _, err := tx.Exec(insertQuery, insertValues...); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
