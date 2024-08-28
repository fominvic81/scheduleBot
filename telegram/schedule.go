package telegram

import (
	"errors"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func GetSchedule(c tele.Context, days int, offset int, startFromMonday bool) ([]api.ScheduleDay, error) {
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

	return schedule, nil
}

func SendSchedule(c tele.Context, formatter func(day *api.ScheduleDay) string, days int, offset int, startFromMonday bool) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	schedule, err := GetSchedule(c, days, offset, startFromMonday)
	if err != nil {
		return err
	}

	if len(schedule) == 0 {
		return c.Send("Розклад не знайдено", c.Get("keyboard").(func() *tele.ReplyMarkup)())
	}

	for _, day := range schedule {
		err = c.Send(formatter(&day), tele.ModeMarkdownV2, c.Get("keyboard").(func() *tele.ReplyMarkup)())

		if err != nil {
			return err
		}
	}

	return nil
}
