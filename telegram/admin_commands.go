package telegram

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func SendHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)
	database := c.Get("database").(*sql.DB)

	if !user.IsAdmin {
		return nil
	}

	if c.Message().ReplyTo == nil {
		if err := c.Send("Команда має бути відповіддю на інше повідомлення"); err != nil {
			return err
		}
		return nil
	}

	payload := c.Message().Payload

	var user2 *db.User = nil

	if payload[0] == '@' {
		var err error
		user2, err = db.GetUserByUsername(database, payload[1:])
		if errors.Is(err, sql.ErrNoRows) {
			if err := c.Send("Користувача не знайдно"); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
	} else {
		id, err := strconv.ParseInt(payload, 10, 64)
		if err != nil {
			if err := c.Send("Не валідний id користувача"); err != nil {
				return err
			}
			return nil
		}
		user2, err = db.GetUser(database, id)
		if errors.Is(err, sql.ErrNoRows) {
			if err := c.Send("Користувача не знайдно"); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
	}

	if _, err := c.Bot().Copy(tele.ChatID(user2.Id), c.Message().ReplyTo); err != nil {
		return err
	}

	return nil
}
