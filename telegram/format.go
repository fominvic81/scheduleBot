package telegram

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func Escape(msg string) string {
	r, _ := regexp.Compile("([\\*\\_\\~\\|\\[\\]\\(\\)\\>\\-\\!\\.\\`])")
	return r.ReplaceAllString(msg, "\\$1")
}

func FormatDay(c tele.Context, day *api.Day) []string {
	user := c.Get("user").(*db.User)

	messages := []string{}

	header := fmt.Sprintf("%s, %s\n\n", Escape(day.WeekDay), Escape(day.Date))
	footer := fmt.Sprintf("Оновлено %s", time.Now().Format("15:04:05"))

	message := header

	for i, class := range day.Classes {
		classMessage := ""
		classMessage += fmt.Sprintf("⚪ *%s*, \\[%s\\-%s\\]\n", Escape(class.StudyTime), Escape(class.Begin), Escape(class.End))
		classMessage += fmt.Sprintf("Предмет: %s\n", Escape(class.Discipline))
		if user.Settings.ShowTeacher {
			classMessage += fmt.Sprintf("Викладач: %s\n", Escape(class.Employee))
		}
		classMessage += fmt.Sprintf("Тип: \\[*%s*\\] Кабінет: \\[*%s*\\]\n", Escape(class.Type), Escape(class.Cabinet))

		if user.Settings.ShowGroups != 0 {
			if class.Groups != "" {
				groups := class.Groups
				if len(groups) > 90 && user.Settings.ShowGroups == 1 {
					count := strings.Count(groups, ",")
					groups = groups[:90]
					count -= strings.Count(groups, ",")

					groups = fmt.Sprintf("%s... (І ще %v)", groups[:strings.LastIndex(groups, ",")], count)
				}
				classMessage += fmt.Sprintf("Групи: %s\n", Escape(groups))
			} else {
				classMessage += Escape("Групи: Пошук...\n")
			}
		}

		classMessage += "\n"

		if len(message)+len(classMessage)+len(footer) >= 4096 || i == len(day.Classes)-1 {
			message += footer
			messages = append(messages, message)

			message = Escape("...\n\n") + header
		}

		message += classMessage
	}

	// messages = append(messages, message)

	return messages
}

func FormatDayShort(_ tele.Context, day *api.Day) []string {
	r, _ := regexp.Compile(`\d*`)

	message := fmt.Sprintf("%s, %s\n\n", Escape(day.WeekDay), Escape(day.Date))

	for _, class := range day.Classes {
		message += fmt.Sprintf("%s: %s, %s\n", r.FindString(class.StudyTime), Escape(class.Discipline), Escape(class.Type))
	}

	return []string{message}
}
