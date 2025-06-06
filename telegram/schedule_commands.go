package telegram

import (
	"fmt"
	"hash/fnv"
	"slices"

	"github.com/fominvic81/scheduleBot/api"

	tele "gopkg.in/telebot.v3"
)

func DayHandler(c tele.Context) error {
	return SendScheduleWithOptions(c, true, FormatDay, 1, 0, false)
}

func NextHandler(c tele.Context) error {
	return SendScheduleWithOptions(c, true, FormatDay, 1, 1, false)
}

func NextNextHandler(c tele.Context) error {
	return SendScheduleWithOptions(c, true, FormatDay, 1, 2, false)
}

func WeekHandler(c tele.Context) error {
	return SendScheduleWithOptions(c, false, FormatDay, 7, 0, true)
}

func NextWeekHandler(c tele.Context) error {
	return SendScheduleWithOptions(c, false, FormatDay, 7, 7, true)
}

func ShortHandler(c tele.Context) error {
	return SendScheduleWithOptions(c, false, FormatDayShort, 21, 0, true)
}

func hash(str string) string {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(str))

	return fmt.Sprintf("%x", hasher.Sum64())
}

func SubjectHandler(c tele.Context) error {
	asked, err := Ask(c)
	if err != nil || asked {
		return err
	}

	start, end := GetDateRange(21, 0, true)
	schedule, err := GetSchedule(c, start, end, true)
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
		rows = append(rows, markup.Row(markup.Data(discipline, "discipline:"+hash(discipline))))
	}
	markup.Inline(rows...)

	return c.Send("Виберіть предмет", &markup)
}

func SendSubject(c tele.Context, discipline string) error {
	return SendScheduleWithOptions(c, true, func(c tele.Context, day *api.Day, withGroups bool) []string {
		day2 := api.Day{
			WeekDay: day.WeekDay,
			Date:    day.Date,
			Classes: make([]api.Class, 0, len(day.Classes)),
		}
		for _, class := range day.Classes {
			if hash(class.Discipline) == discipline {
				day2.Classes = append(day2.Classes, class)
			}
		}
		return FormatDay(c, &day2, withGroups)
	}, 21, 0, true)
}
