package telegram

import (
	"errors"

	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func SetGroup(c tele.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return errors.New("failed to get user in 'SetGroup'")
	}

	user.StudyGroup = nil

	err := user.Save()
	if err != nil {
		return err
	}

	_, err = Ask(c)
	if err != nil {
		return err
	}

	return nil
}
