package telegram

import (
	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func TextHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)

	keyboards := GetKeyboards()
	for _, keyboard := range keyboards {
		for _, row := range keyboard {
			for _, btn := range row {
				if c.Text() == btn.Text {
					if user.State != db.UserStateNone {
						user.State = db.UserStateNone
						if err := user.Save(); err != nil {
							return err
						}
					}
					if err := btn.Handler(c); err != nil {
						return err
					}
					return nil
				}
			}
		}
	}

	switch user.State {
	case db.UserStateSearchGroup:
		if err := SearchGroupHandler(c); err != nil {
			return err
		}
	case db.UserStateSearchTeacher:
		if err := SearchTeacherHandler(c); err != nil {
			return err
		}
	}

	return nil
}
