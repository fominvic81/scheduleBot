package telegram

import (
	"errors"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func GetSchedule(c tele.Context, withGroups bool, days int, offset int, startFromMonday bool) ([]api.Day, error) {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return nil, errors.New("failed to get user in 'GetSchedule'")
	}

	if user.StudyGroup == nil {
		return nil, errors.New("failed to get schedule, used does not have selected study group")
	}

	now := time.Now()

	startOffset := 0
	if startFromMonday {
		startOffset = -int(now.Weekday())
	}

	start := now.AddDate(0, 0, startOffset+offset)
	end := start.AddDate(0, 0, days-1)

	if now.After(start) {
		start = now
		if start.After(end) {
			end = start
		}
	}

	schedule, err := api.GetSchedule(*user.StudyGroup, start, end)
	if err != nil {
		return nil, err
	}

	if withGroups {
		err = api.GetScheduleGroups(schedule, start, end)
		if err != nil {
			LogError(err, c)
		}
	}

	return schedule, nil
}

func GetDayMarkup(c tele.Context, date string) *tele.ReplyMarkup {
	time, err := time.Parse("02.01.2006", date)

	if err != nil {
		LogError(err, c)
		return &tele.ReplyMarkup{}
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Попередній день", Data: "update:" + time.AddDate(0, 0, -1).Format("02.01.2006")},
				{Text: "Наступний день", Data: "update:" + time.AddDate(0, 0, 1).Format("02.01.2006")},
			},
			{{Text: "Оновити", Data: "update:" + date}},
		},
	}
}

func SendSchedule(c tele.Context, withGroups bool, formatter func(day *api.Day) string, days int, offset int, startFromMonday bool) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	schedule, err := GetSchedule(c, withGroups, days, offset, startFromMonday)
	if err != nil {
		return err
	}

	if len(schedule) == 0 {
		return c.Send("Розклад не знайдено", GetMarkup(c, nil))
	}

	for _, day := range schedule {
		if days > 1 {
			err = c.Send(formatter(&day), tele.ModeMarkdownV2, GetMarkup(c, nil))
		} else {
			err = c.Send(formatter(&day), tele.ModeMarkdownV2, GetMarkup(c, GetDayMarkup(c, day.Date)))
		}

		if err != nil {
			return err
		}
	}

	return nil
}
