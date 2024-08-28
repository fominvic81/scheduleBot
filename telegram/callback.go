package telegram

import (
	"errors"
	"regexp"

	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func CallbackData(c tele.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return errors.New("failed to get user in CallbackData")
	}

	r, _ := regexp.Compile("([a-z]+):(.*)")
	matches := r.FindStringSubmatch(c.Data())

	if len(matches) == 3 {
		key := matches[1]
		value := matches[2]

		switch key {
		case "faculty":
			user.Faculty = &value
			_, err := Ask(c)
			if err != nil {
				return err
			}
		case "form":
			user.EducationForm = &value
			_, err := Ask(c)
			if err != nil {
				return err
			}
		case "course":
			user.Course = &value
			_, err := Ask(c)
			if err != nil {
				return err
			}
		case "group":
			user.StudyGroup = &value

			err := c.Send("Готово!")
			if err != nil {
				return err
			}
		case "discipline":
			err := SendSubject(c, value)
			if err != nil {
				return err
			}
		}

		msg := c.Message()
		if msg != nil {
			c.Bot().Delete(msg)
		}

		user.Save()
	}

	return nil
}
