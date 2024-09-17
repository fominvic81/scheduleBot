package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func boolToInt(a bool) int {
	if a {
		return 1
	}
	return 0
}

func CallbackData(c tele.Context) error {
	user := c.Get("user").(*db.User)

	r, _ := regexp.Compile("([a-z/]+)(?::([^;]+);?(.*))?")
	matches := r.FindStringSubmatch(c.Data())

	if len(matches) >= 3 {
		key := matches[1]
		value := matches[2]

		switch key {
		case "faculty":
			user.Faculty = &value
			_, err := Ask(c)
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
		case "form":
			user.EducationForm = &value
			_, err := Ask(c)
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
		case "course":
			user.Course = &value
			_, err := Ask(c)
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
		case "group":
			user.StudyGroup = &value
			err := c.Send("Готово!", GetMarkup(c, nil))
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
		case "discipline":
			err := SendSubject(c, value)
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
		case "update":
			if user.StudyGroup == nil {
				_, err := Ask(c)
				return err
			}
			date, err := time.Parse("02.01.2006", value)
			if err != nil {
				return err
			}

			err = SendSchedule(c, c.Message(), true, FormatDay, date, date)
			if err != nil {
				return err
			}
		case "settings":
			err := c.Edit("Налаштування", GetMarkup(c, &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{tele.InlineButton{Text: "Формат розкладу", Data: "settings/format"}},
					{tele.InlineButton{Text: "Вибіркові дисципліни", Data: "settings/disciplines"}},
					{tele.InlineButton{Text: "Закрити", Data: "delete"}},
				},
			}))
			if err != nil {
				return err
			}
		case "settings/format":
			switch value {
			case "show-groups":
				user.Settings.ShowGroups = (user.Settings.ShowGroups + 1) % 3
			case "show-teacher":
				user.Settings.ShowTeacher = !user.Settings.ShowTeacher
			}
			err := user.Save()
			if err != nil {
				return err
			}

			showGroups := []string{"Не показувати", "Частково", "Всі"}[user.Settings.ShowGroups]
			showTeacher := []string{"Не показувати", "Показувати"}[boolToInt(user.Settings.ShowTeacher)]
			err = c.Edit("Формат повідомлення", GetMarkup(c, &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{{Text: "Групи: " + showGroups, Data: "settings/format:show-groups"}},
					{{Text: "Викладач: " + showTeacher, Data: "settings/format:show-teacher"}},
					{{Text: "Назад", Data: "settings"}},
				},
			}))
			if err != nil {
				return err
			}
		case "settings/disciplines":
			schedule, err := GetSchedule(c, time.Now().AddDate(0, 0, -14), time.Now().AddDate(0, 0, 14), false)
			if err != nil {
				return err
			}

			subjects := map[string]string{}
			for _, day := range schedule {
				for _, class := range day.Classes {
					subjects[hash(class.Discipline)] = class.Discipline
				}
			}

			if value == "all-off" {
				user.Settings.HiddenSubjects = []string{}
				for _, subject := range subjects {
					user.Settings.HiddenSubjects = append(user.Settings.HiddenSubjects, subject)
				}
				err = user.Save()
				if err != nil {
					return err
				}
			}
			if value == "all-on" {
				user.Settings.HiddenSubjects = []string{}
				err = user.Save()
				if err != nil {
					return err
				}
			}
			if len(matches) >= 4 {
				status := matches[3] == "on"
				fmt.Println(status, matches[3])
				if clicked, ok := subjects[value]; ok {
					current := !slices.Contains(user.Settings.HiddenSubjects, clicked)
					if current != status {
						if status {
							var newSubjects []string
							for _, subject := range user.Settings.HiddenSubjects {
								if subject != clicked {
									newSubjects = append(newSubjects, subject)
								}
							}
							user.Settings.HiddenSubjects = newSubjects
						} else {
							user.Settings.HiddenSubjects = append(user.Settings.HiddenSubjects, clicked)
						}
						err = user.Save()
						if err != nil {
							return err
						}
					}
				}
			}

			var keyboard [][]tele.InlineButton
			for hashed, subject := range subjects {
				hidden := slices.Contains(user.Settings.HiddenSubjects, subject)
				prefix := []string{"❌", "✅"}[boolToInt(!hidden)]
				statusText := "off"
				if hidden {
					statusText = "on"
				}

				keyboard = append(keyboard, []tele.InlineButton{{
					Text: fmt.Sprintf("%s %s", prefix, subject),
					Data: "settings/disciplines:" + hashed + ";" + statusText,
				}})
			}
			slices.SortFunc(keyboard, func(a []tele.InlineButton, b []tele.InlineButton) int {
				// ❌ and ✅ take 3 bytes each, so start sorting from fourth byte
				return strings.Compare(a[0].Text[3:], b[0].Text[3:])
			})
			keyboard = append(keyboard, []tele.InlineButton{
				{Text: "Вимкнути всі ❌", Data: "settings/disciplines:all-off"},
				{Text: "Ввімкнути всі ✅", Data: "settings/disciplines:all-on"},
			})
			keyboard = append(keyboard, []tele.InlineButton{{Text: "Назад ↩️", Data: "settings"}})

			err = c.Edit("Вибіркові дисципліни", GetMarkup(c, &tele.ReplyMarkup{
				InlineKeyboard: keyboard,
			}))
			if errors.Is(err, tele.ErrSameMessageContent) {
				err = nil
			}
		case "delete":
			err := c.Delete()
			if err != nil {
				return err
			}
		}

		err := user.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
