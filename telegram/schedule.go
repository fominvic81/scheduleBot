package telegram

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func GetDateRange(days int, offset int, startFromMonday bool) (time.Time, time.Time) {
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

	return start, end
}

func GetSchedule(c tele.Context, start time.Time, end time.Time, filter bool) ([]api.Day, error) {
	user := c.Get("user").(*db.User)

	if user.StudyGroup == nil {
		return nil, errors.New("failed to get schedule, used does not have selected study group")
	}

	unfilteredSchedule, err := api.GetSchedule(*user.StudyGroup, start, end)
	if err != nil {
		return nil, err
	}

	schedule := []api.Day{}
	for _, day := range unfilteredSchedule {
		classes := []api.Class{}
		for _, class := range day.Classes {
			if !filter || !slices.Contains(user.Settings.HiddenSubjects, class.Discipline) {
				classes = append(classes, class)
			}
		}
		day.Classes = classes
		schedule = append(schedule, day)
	}

	return schedule, nil
}

func GetScheduleGroups(c tele.Context, schedule []api.Day, start time.Time, end time.Time) {
	err := api.GetScheduleGroups(schedule, start, end)
	if err != nil {
		LogError(err, c)
	}
}

func GetDayMarkup(c tele.Context, date string) *tele.ReplyMarkup {
	current, err := time.Parse("02.01.2006", date)

	if err != nil {
		LogError(err, c)
		return &tele.ReplyMarkup{}
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "⏪", Data: "update:" + current.AddDate(0, 0, -1).Format("02.01.2006") + ";prev"},
				{Text: "Сьогодні", Data: "update:" + time.Now().Format("02.01.2006") + ";today"},
				{Text: "⏩", Data: "update:" + current.AddDate(0, 0, 1).Format("02.01.2006") + ";next"},
			},
			{{Text: "Оновити", Data: "update:" + date}},
		},
	}
}

func SendSchedule(c tele.Context, withGroups bool, formatter func(c tele.Context, day *api.Day) string, days int, offset int, startFromMonday bool) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	start, end := GetDateRange(days, offset, startFromMonday)
	schedule, err := GetSchedule(c, start, end, true)
	if err != nil {
		return err
	}

	if len(schedule) == 0 {
		return c.Send("Розклад не знайдено", GetMarkup(c, nil))
	}

	for i, day := range schedule {
		if days > 1 {
			msg, err := c.Bot().Send(c.Recipient(), formatter(c, &day), tele.ModeMarkdownV2, GetMarkup(c, nil))
			if err != nil {
				return err
			}
			schedule[i].MessageId = msg.ID
		} else {
			msg, err := c.Bot().Send(c.Recipient(), formatter(c, &day), tele.ModeMarkdownV2, GetMarkup(c, GetDayMarkup(c, day.Date)))
			if err != nil {
				return err
			}
			schedule[i].MessageId = msg.ID
		}
	}

	if withGroups {
		GetScheduleGroups(c, schedule, start, end)

		for _, day := range schedule {
			msg := tele.StoredMessage{
				ChatID:    c.Chat().ID,
				MessageID: fmt.Sprintf("%v", day.MessageId),
			}
			if days > 1 {
				_, err := c.Bot().Edit(msg, formatter(c, &day), tele.ModeMarkdownV2, GetMarkup(c, nil))
				if err != nil && err != tele.ErrSameMessageContent {
					return err
				}
			} else {
				_, err := c.Bot().Edit(msg, formatter(c, &day), tele.ModeMarkdownV2, GetMarkup(c, GetDayMarkup(c, day.Date)))
				if err != nil && err != tele.ErrSameMessageContent {
					return err
				}
			}
		}
	}

	return nil
}
