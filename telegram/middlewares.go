package telegram

import (
	"database/sql"
	"errors"

	"github.com/fominvic81/scheduleBot/consts"
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
		if sender == nil {
			return errors.New("failed to identify user")
		}
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
		return next(c)
	}
}

func LogMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		LogAction(c)
		return next(c)
	}
}

func KeyboardMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user := c.Get("user").(*db.User)

		keyboardUsed := false
		if user.KeyboardVersion != consts.KeyboardVersion {
			c.Set("keyboard", func() [][]tele.ReplyButton {
				keyboardUsed = true
				return GetReplyKeyboard()
			})
		} else {
			c.Set("keyboard", func() [][]tele.ReplyButton {
				return nil
			})
		}

		err := next(c)
		if err != nil {
			return err
		}

		if keyboardUsed {
			user.KeyboardVersion = consts.KeyboardVersion
		}

		return user.Save()
	}
}

func GetMarkup(c tele.Context, markup *tele.ReplyMarkup) *tele.ReplyMarkup {
	if markup == nil {
		markup = &tele.ReplyMarkup{}
	}
	if markup.ReplyKeyboard == nil {
		markup.ReplyKeyboard = c.Get("keyboard").(func() [][]tele.ReplyButton)()
		markup.ResizeKeyboard = true
	}
	return markup
}
