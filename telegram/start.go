package telegram

import (
	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func StartHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)

	if err := c.Send(HelpMsg()); err != nil {
		return err
	}

	if user.StudyGroup == nil {
		if err := AskGroup(c); err != nil {
			return err
		}
		return nil
	}

	return nil
}
