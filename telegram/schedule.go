package telegram

import (
	"errors"
	"fmt"
	"github.com/fominvic81/scheduleBot/consts"
	"slices"
	"sync"
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

	var schedule []api.Day
	for _, day := range unfilteredSchedule {
		var classes []api.Class
		for _, class := range day.Classes {
			if !filter || !slices.Contains(user.Settings.HiddenSubjects, class.Discipline) {
				classes = append(classes, class)
			}
		}
		day.Classes = classes
		if len(classes) > 0 {
			schedule = append(schedule, day)
		}
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

var msgLastModMu = sync.Mutex{}
var msgLastMod = map[string]int{}

func SendSchedule(c tele.Context, message *tele.Message, withGroups bool, formatter func(c tele.Context, day *api.Day) string, start time.Time, end time.Time) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	schedule, err := GetSchedule(c, start, end, true)
	if err != nil {
		return err
	}

	if len(schedule) == 0 {
		var markup *tele.ReplyMarkup
		if start.Day() == end.Day() {
			markup = GetDayMarkup(c, start.Format("02.01.2006"))
		}
		text := consts.WeekDays[start.Weekday()] + ", " + start.Format("02.01.2006")
		if start.Day() != end.Day() {
			text += " - " + consts.WeekDays[end.Weekday()] + ", " + end.Format("02.01.2006")
		}

		text += "\n\nРозклад порожній"
		if message != nil {
			_, err = c.Bot().Edit(message, Escape(text), tele.ModeMarkdownV2, GetMarkup(c, markup))

			if errors.Is(err, tele.ErrSameMessageContent) {
				err = nil
			}
			return err
		} else {
			return c.Send(Escape(text), tele.ModeMarkdownV2, GetMarkup(c, markup))
		}
	}

	for i, day := range schedule {
		text := formatter(c, &day)
		var markup *tele.ReplyMarkup
		if start.Day() == end.Day() {
			markup = GetDayMarkup(c, day.Date)
		}

		var msg *tele.Message
		if message != nil && start.Day() == end.Day() {
			msg, err = c.Bot().Edit(message, text, tele.ModeMarkdownV2, GetMarkup(c, markup))
		} else {
			msg, err = c.Bot().Send(c.Recipient(), text, tele.ModeMarkdownV2, GetMarkup(c, markup))
		}
		if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
			return err
		}
		schedule[i].MessageId = msg.ID
	}

	if withGroups {
		key := fmt.Sprintf("%v|%v", c.Chat().ID, c.Message().ID)
		msgLastModMu.Lock()
		count := msgLastMod[key] + 1
		msgLastMod[key] = count
		msgLastModMu.Unlock()

		GetScheduleGroups(c, schedule, start, end)

		msgLastModMu.Lock()
		current := msgLastMod[key]
		msgLastModMu.Unlock()

		if current == count {
			for _, day := range schedule {
				msg := tele.StoredMessage{
					ChatID:    c.Chat().ID,
					MessageID: fmt.Sprintf("%v", day.MessageId),
				}
				text := formatter(c, &day)
				var markup *tele.ReplyMarkup
				if start.Day() == end.Day() {
					markup = GetDayMarkup(c, day.Date)
				}

				_, err = c.Bot().Edit(&msg, text, tele.ModeMarkdownV2, GetMarkup(c, markup))
				if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
					return err
				}
			}
		}
	}

	return nil
}

func SendScheduleWithOptions(c tele.Context, withGroups bool, formatter func(c tele.Context, day *api.Day) string, days int, offset int, startFromMonday bool) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	start, end := GetDateRange(days, offset, startFromMonday)

	return SendSchedule(c, nil, withGroups, formatter, start, end)
}
