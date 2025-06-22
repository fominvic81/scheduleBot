package telegram

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fominvic81/scheduleBot/consts"

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
		LogError(c, err)
	}
}

func GetDayMarkup(c tele.Context, date string, messageIdsToDelete []string) *tele.ReplyMarkup {
	current, err := time.Parse("02.01.2006", date)

	if err != nil {
		LogError(c, err)
		return &tele.ReplyMarkup{}
	}

	msgsToDelete := ""
	if len(messageIdsToDelete) > 0 {
		msgsToDelete = "|" + strings.Join(messageIdsToDelete, ",")
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "⏪", Data: "update:" + current.AddDate(0, 0, -1).Format("02.01.2006") + msgsToDelete + ";prev"},
				{Text: "Сьогодні", Data: "update:" + "today" + msgsToDelete + ";today"},
				{Text: "⏩", Data: "update:" + current.AddDate(0, 0, 1).Format("02.01.2006") + msgsToDelete + ";next"},
			},
			{{Text: "Оновити 🔄", Data: "update:" + date + msgsToDelete}},
		},
	}
}

func GetTeacherDayMarkup(c tele.Context, employee string, date string, messageIdsToDelete []string) *tele.ReplyMarkup {
	current, err := time.Parse("02.01.2006", date)

	if err != nil {
		LogError(c, err)
		return &tele.ReplyMarkup{}
	}

	msgsToDelete := ""
	if len(messageIdsToDelete) > 0 {
		msgsToDelete = "|" + strings.Join(messageIdsToDelete, ",")
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "⏪", Data: "update-teacher:" + employee + "|" + current.AddDate(0, 0, -1).Format("02.01.2006") + msgsToDelete + ";prev"},
				{Text: "Сьогодні", Data: "update-teacher:" + employee + "|" + "today" + msgsToDelete + ";today"},
				{Text: "⏩", Data: "update-teacher:" + employee + "|" + current.AddDate(0, 0, 1).Format("02.01.2006") + msgsToDelete + ";next"},
			},
			{{Text: "Оновити 🔄", Data: "update-teacher:" + employee + "|" + date + msgsToDelete}},
		},
	}
}

var msgLastModMu = sync.Mutex{}
var msgLastMod = map[string]int{}

func SendSchedule(c tele.Context, message *tele.Message, withGroups bool, formatter Formatter, start time.Time, end time.Time) error {
	user := c.Get("user").(*db.User)
	if user.StudyGroup == nil {
		if err := AskGroup(c); err != nil {
			return err
		}
	}

	schedule, err := GetSchedule(c, start, end, true)
	if err != nil {
		return err
	}

	if len(schedule) == 0 {
		var markup *tele.ReplyMarkup
		if start.Day() == end.Day() {
			markup = GetDayMarkup(c, start.Format("02.01.2006"), nil)
		}
		text := consts.WeekDays[start.Weekday()] + ", " + start.Format("02.01.2006")
		if start.Day() != end.Day() {
			text += " - " + consts.WeekDays[end.Weekday()] + ", " + end.Format("02.01.2006")
		}

		text += "\n\nРозклад порожній"
		if message != nil {
			_, err = c.Bot().Edit(message, Escape(text), tele.ModeMarkdownV2, markup)

			if errors.Is(err, tele.ErrSameMessageContent) {
				err = nil
			}
			return err
		} else {
			return c.Send(Escape(text), tele.ModeMarkdownV2, markup)
		}
	}

	for i, day := range schedule {
		texts := formatter(c, &day, withGroups)

		messageIdsToDelete := []string{}
		var msg *tele.Message
		for j, text := range texts {
			lastIteration := j == len(texts)-1

			var markup *tele.ReplyMarkup
			if start.Day() == end.Day() && lastIteration {
				markup = GetDayMarkup(c, day.Date, messageIdsToDelete)
			}

			if message != nil && start.Day() == end.Day() {
				msg, err = c.Bot().Edit(message, text, tele.ModeMarkdownV2, markup)
			} else {
				msg, err = c.Bot().Send(c.Recipient(), text, tele.ModeMarkdownV2, markup)
			}
			if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
				return err
			}
			schedule[i].MessageIds = append(schedule[i].MessageIds, msg.ID)
			if !lastIteration {
				messageIdsToDelete = append(messageIdsToDelete, fmt.Sprintf("%d", msg.ID))
			}
		}
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
				texts := formatter(c, &day, withGroups)
				messageIdsToDelete := []string{}

				for i, text := range texts {
					lastIteration := i == len(texts)-1

					var markup *tele.ReplyMarkup
					if start.Day() == end.Day() && lastIteration {
						markup = GetDayMarkup(c, day.Date, messageIdsToDelete)
					}

					if i < len(day.MessageIds) {
						msg := tele.StoredMessage{
							ChatID:    c.Chat().ID,
							MessageID: fmt.Sprintf("%d", day.MessageIds[i]),
						}

						_, err = c.Bot().Edit(&msg, text, tele.ModeMarkdownV2, markup)
						if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
							return err
						}
						if !lastIteration {
							messageIdsToDelete = append(messageIdsToDelete, msg.MessageID)
						}
					} else {
						msg, err := c.Bot().Send(c.Chat(), text, tele.ModeMarkdownV2, markup)
						if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
							return err
						}
						messageIdsToDelete = append(messageIdsToDelete, fmt.Sprintf("%d", msg.ID))
					}
				}
			}
		}
	}

	return nil
}

func SendScheduleWithOptions(c tele.Context, withGroups bool, formatter Formatter, days int, offset int, startFromMonday bool) error {
	user := c.Get("user").(*db.User)
	if user.StudyGroup == nil {
		if err := AskGroup(c); err != nil {
			return err
		}
		return nil
	}

	start, end := GetDateRange(days, offset, startFromMonday)

	return SendSchedule(c, nil, withGroups, formatter, start, end)
}
