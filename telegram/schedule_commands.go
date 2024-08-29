package telegram

import (
	"slices"

	"github.com/fominvic81/scheduleBot/api"

	tele "gopkg.in/telebot.v3"
)

func Day(c tele.Context) error {
	return SendSchedule(c, true, FormatDay, 1, 0, false)
}

func Next(c tele.Context) error {
	return SendSchedule(c, true, FormatDay, 1, 1, false)
}

func NextNext(c tele.Context) error {
	return SendSchedule(c, true, FormatDay, 1, 2, false)
}

func Week(c tele.Context) error {
	return SendSchedule(c, true, FormatDay, 7, 0, true)
}

func NextWeek(c tele.Context) error {
	return SendSchedule(c, true, FormatDay, 7, 7, true)
}

func Short(c tele.Context) error {
	return SendSchedule(c, true, FormatDayShort, 21, 0, true)
}

func Subject(c tele.Context) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	schedule, err := GetSchedule(c, false, 21, 0, true)
	if err != nil {
		return err
	}

	if len(schedule) == 0 {
		return c.Send("Розклад не знайдено")
	}

	disciplines := make([]string, 0, 6)

	for _, day := range schedule {
		for _, class := range day.Classes {
			if !slices.Contains(disciplines, class.Discipline) {
				disciplines = append(disciplines, class.Discipline)
			}
		}
	}

	slices.Sort(disciplines)

	markup := tele.ReplyMarkup{}
	rows := make([]tele.Row, 0, 6)
	for _, discipline := range disciplines {
		rows = append(rows, markup.Row(markup.Data(discipline, "discipline:"+discipline)))
	}
	markup.Inline(rows...)

	return c.Send("Виберіть предмет", &markup)
}

func SendSubject(c tele.Context, discipline string) error {
	return SendSchedule(c, true, func(day *api.Day) string {
		day2 := api.Day{
			WeekDay: day.WeekDay,
			Date:    day.Date,
			Classes: make([]api.Class, 0, len(day.Classes)),
		}
		for _, class := range day.Classes {
			if class.Discipline == discipline {
				day2.Classes = append(day2.Classes, class)
			}
		}
		return FormatDay(&day2)
	}, 21, 0, true)
}
