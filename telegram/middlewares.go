package telegram

import (
	"database/sql"

	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func EmptyStrAsNil(str string) *string {
	if len(str) == 0 {
		return nil
	}
	return &str
}

func UserMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		database := c.Get("database").(*sql.DB)

		sender := c.Sender()
		if sender != nil {
			user, err := db.GetOrCreateUser(database, sender.ID, sender.FirstName)
			if err != nil {
				return err
			}
			user.Firstname = sender.FirstName
			user.Lastname = EmptyStrAsNil(sender.LastName)
			user.Username = EmptyStrAsNil(sender.Username)
			user.Messages += 1
			err = user.Save()

			if err != nil {
				return err
			}

			c.Set("user", user)
		}
		return next(c)
	}
}

func KeyboardMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user, ok := c.Get("user").(*db.User)
		if !ok {
			return next(c)
		}

		keyboardUsed := false
		if user.KeyboardVersion != db.KeyboardVersion {
			c.Set("keyboard", func() *tele.ReplyMarkup {
				keyboardUsed = true
				return GetKeyboardMarkup()
			})
		} else {
			c.Set("keyboard", func() *tele.ReplyMarkup {
				return nil
			})
		}

		err := next(c)
		if err != nil {
			return err
		}

		if keyboardUsed {
			user.KeyboardVersion = db.KeyboardVersion
		}

		return user.Save()
	}
}
