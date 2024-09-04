package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/consts"
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

		var err error = nil
		delete := true

		switch key {
		case "faculty":
			user.Faculty = &value
			_, err = Ask(c)
		case "form":
			user.EducationForm = &value
			_, err = Ask(c)
		case "course":
			user.Course = &value
			_, err = Ask(c)
		case "group":
			user.StudyGroup = &value
			err = c.Send("Готово!", GetMarkup(c, nil))
		case "discipline":
			err = SendSubject(c, value)
		case "update":
			delete = false
			if user.StudyGroup == nil {
				_, err = Ask(c)
				return err
			}
			var date time.Time
			date, err = time.Parse("02.01.2006", value)
			if err != nil {
				return err
			}

			var days []api.Day
			days, err = api.GetSchedule(*user.StudyGroup, date, date)
			if err != nil {
				return err
			}

			err = api.GetScheduleGroup(days, date, date)
			if err != nil {
				return err
			}

			markup := GetDayMarkup(date.Format("02.01.2006"))
			if len(days) == 0 {
				err = c.Edit(consts.WeekDays[int(date.Weekday())]+", "+date.Format("02.01.2006")+"\n\nРозкладу немає", markup)
			} else {
				text := FormatDay(&days[0]) + fmt.Sprintf("Оновлено %s", time.Now().Format("15:04:05"))
				err = c.Edit(text, tele.ModeMarkdownV2, markup)
			}

			if err == tele.ErrSameMessageContent {
				err = nil
			}
		}

		if err != nil {
			return err
		}

		if delete {
			err = c.Delete()
			if err != nil {
				return err
			}
		}

		err = user.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
