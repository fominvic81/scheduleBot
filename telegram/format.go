package telegram

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func Escape(msg string) string {
	r, _ := regexp.Compile(`([\.\-\_\*\(\)\!])`)

	return r.ReplaceAllString(msg, "\\$1")
}

func FormatDay(c tele.Context, day *api.Day) string {
	user := c.Get("user").(*db.User)

	message := fmt.Sprintf("%s, %s\n\n", Escape(day.WeekDay), Escape(day.Date))

	for _, class := range day.Classes {
		message += fmt.Sprintf("⚪ *%s*, \\[%s\\-%s\\]\n", Escape(class.StudyTime), Escape(class.Begin), Escape(class.End))
		message += fmt.Sprintf("Предмет: %s\n", Escape(class.Discipline))
		if user.Settings.ShowTeacher {
			message += fmt.Sprintf("Викладач: %s\n", Escape(class.Employee))
		}
		message += fmt.Sprintf("Тип: \\[*%s*\\] Кабінет: \\[*%s*\\]\n", Escape(class.Type), Escape(class.Cabinet))

		if user.Settings.ShowGroups != 0 {
			if class.Groups != "" {
				groups := class.Groups
				if len(groups) > 90 && user.Settings.ShowGroups == 1 {
					count := strings.Count(groups, ",")
					groups = groups[:90]
					count -= strings.Count(groups, ",")

					groups = fmt.Sprintf("%s... (І ще %v)", groups[:strings.LastIndex(groups, ",")], count)
				}
				message += fmt.Sprintf("Групи: %s\n", Escape(groups))
			} else {
				message += Escape("Групи: Пошук...\n")
			}
		}
		message += "\n"
	}

	return message
}

func FormatDayShort(c tele.Context, day *api.Day) string {
	r, _ := regexp.Compile(`\d*`)

	message := fmt.Sprintf("%s, %s\n\n", Escape(day.WeekDay), Escape(day.Date))

	for _, class := range day.Classes {
		message += fmt.Sprintf("%s: %s, %s\n", r.FindString(class.StudyTime), Escape(class.Discipline), Escape(class.Type))
	}

	return message
}
