package telegram

import (
	"fmt"
	"regexp"

	"github.com/fominvic81/scheduleBot/api"
)

func Escape(msg string) string {
	r, _ := regexp.Compile(`([\.\-\_\*\(\)\!])`)

	return r.ReplaceAllString(msg, "\\$1")
}

func FormatDay(day *api.Day) string {
	message := fmt.Sprintf("%s, %s\n\n", Escape(day.WeekDay), Escape(day.Date))

	for _, class := range day.Classes {
		message += fmt.Sprintf("⚪ *%s*, \\[%s\\-%s\\]\n", Escape(class.StudyTime), Escape(class.Begin), Escape(class.End))
		message += fmt.Sprintf("Предмет: %s\n", Escape(class.Discipline))
		message += fmt.Sprintf("Викладач: %s\n", Escape(class.Employee))
		message += fmt.Sprintf("Тип: \\[*%s*\\] Кабінет: \\[*%s*\\]\n", Escape(class.Type), Escape(class.Cabinet))

		if class.Groups != "" {
			message += fmt.Sprintf("Групи: %s\n", Escape(class.Groups))
		}
		message += "\n"
	}

	return message
}

func FormatDayShort(day *api.Day) string {
	r, _ := regexp.Compile(`\d*`)

	message := fmt.Sprintf("%s, %s\n\n", Escape(day.WeekDay), Escape(day.Date))

	for _, class := range day.Classes {
		message += fmt.Sprintf("%s: %s, %s\n", r.FindString(class.StudyTime), Escape(class.Discipline), Escape(class.Type))
	}

	return message
}
