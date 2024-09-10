package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/consts"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

var msgLastModMu = sync.Mutex{}
var msgLastMod = map[string]int{}

func boolToInt(a bool) int {
	if a {
		return 1
	}
	return 0
}

func CallbackData(c tele.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return errors.New("failed to get user in CallbackData")
	}

	r, _ := regexp.Compile("([a-z/]+)(?::([^;]+);?(.*))?")
	matches := r.FindStringSubmatch(c.Data())

	if len(matches) >= 3 {
		key := matches[1]
		value := matches[2]

		var err error = nil
		delete := false

		switch key {
		case "faculty":
			user.Faculty = &value
			_, err = Ask(c)
			delete = true
		case "form":
			user.EducationForm = &value
			_, err = Ask(c)
			delete = true
		case "course":
			user.Course = &value
			_, err = Ask(c)
			delete = true
		case "group":
			user.StudyGroup = &value
			err = c.Send("Готово!", GetMarkup(c, nil))
			delete = true
		case "discipline":
			delete = true
			err = SendSubject(c, value)
		case "update":
			if user.StudyGroup == nil {
				_, err = Ask(c)
				return err
			}
			var date time.Time
			date, err = time.Parse("02.01.2006", value)
			if err != nil {
				return err
			}

			var days []api.Day
			days, err = GetSchedule(c, date, date, true)
			if err != nil {
				return err
			}

			markup := GetDayMarkup(c, date.Format("02.01.2006"))
			if len(days) == 0 {
				err = c.Edit(consts.WeekDays[int(date.Weekday())]+", "+date.Format("02.01.2006")+"\n\nРозкладу немає", markup)
			} else {
				text := FormatDay(c, &days[0]) + fmt.Sprintf("Оновлено %s", time.Now().Format("15:04:05"))
				err = c.Edit(text, tele.ModeMarkdownV2, markup)

				if err != nil && err != tele.ErrSameMessageContent {
					return err
				}
				key := fmt.Sprintf("%v|%v", c.Chat().ID, c.Message().ID)

				msgLastModMu.Lock()
				count := msgLastMod[key] + 1
				msgLastMod[key] = count
				msgLastModMu.Unlock()

				GetScheduleGroups(c, days, date, date)

				msgLastModMu.Lock()
				current := msgLastMod[key]
				msgLastModMu.Unlock()

				if current == count {
					if err != nil {
						LogError(err, c)
					} else {
						text = FormatDay(c, &days[0]) + fmt.Sprintf("Оновлено %s", time.Now().Format("15:04:05"))
						err = c.Edit(text, tele.ModeMarkdownV2, markup)
					}
				}
			}

			if err == tele.ErrSameMessageContent {
				err = nil
			}
		case "settings":
			err = c.Edit("Налаштування", GetMarkup(c, &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{tele.InlineButton{Text: "Формат розкладу", Data: "settings/format"}},
					{tele.InlineButton{Text: "Вибіркові дисципліни", Data: "settings/disciplines"}},
					{tele.InlineButton{Text: "Закрити", Data: "delete"}},
				},
			}))
		case "settings/format":
			switch value {
			case "show-groups":
				user.Settings.ShowGroups = (user.Settings.ShowGroups + 1) % 3
			case "show-teacher":
				user.Settings.ShowTeacher = !user.Settings.ShowTeacher
			}
			err = user.Save()
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
		case "settings/disciplines":
			var schedule []api.Day
			schedule, err = GetSchedule(c, time.Now().AddDate(0, 0, -14), time.Now().AddDate(0, 0, 14), false)
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
				if clicked, ok := subjects[value]; ok {
					current := slices.Contains(user.Settings.HiddenSubjects, clicked)
					if current != status {
						if status {
							user.Settings.HiddenSubjects = append(user.Settings.HiddenSubjects, clicked)
						} else {
							newSubjects := []string{}
							for _, subject := range user.Settings.HiddenSubjects {
								if subject != clicked {
									newSubjects = append(newSubjects, subject)
								}
							}
							user.Settings.HiddenSubjects = newSubjects
						}
						err = user.Save()
						if err != nil {
							return err
						}
					}
				}
			}

			keyboard := [][]tele.InlineButton{}
			for hashed, subject := range subjects {
				hidden := slices.Contains(user.Settings.HiddenSubjects, subject)
				prefix := []string{"❌", "✅"}[boolToInt(!hidden)]
				statusText := "on"
				if hidden {
					statusText = "off"
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
			if err == tele.ErrSameMessageContent {
				err = nil
			}
		case "delete":
			delete = true
		}

		if err != nil {
			return err
		}

		if delete {
			err = c.Delete()
			if err != nil {
				return err
			}
		}

		err = user.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
