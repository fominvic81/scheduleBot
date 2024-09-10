package telegram

import (
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func SetGroup(c tele.Context) error {
	user := c.Get("user").(*db.User)

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
