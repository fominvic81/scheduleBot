package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/consts"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func boolToInt(a bool) int {
	if a {
		return 1
	}
	return 0
}

func DeleteMessage(c tele.Context) error {
	err := c.Delete()
	if errors.Is(err, tele.ErrNotFoundToDelete) {
		LogError(c, err)
	} else if errors.Is(err, tele.ErrNoRightsToDelete) {
		if err := c.Edit("Видалено"); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func CallbackDataHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)

	r := regexp.MustCompile(`([a-z/\-]+)(?::([^;]+))?(?:;.+)?`)
	matches := r.FindStringSubmatch(c.Data())

	if len(matches) >= 3 {
		key := matches[1]
		values := strings.Split(matches[2], "|")
		value := values[0]
		value2 := ""
		if len(values) > 1 {
			value2 = values[1]
		}
		value3 := ""
		if len(values) > 2 {
			value3 = values[2]
		}

		switch key {
		case "ask-faculty":
			if err := AskFaculty(c); err != nil {
				return err
			}
		case "ask-form":
			if err := AskForm(c); err != nil {
				return err
			}
		case "ask-course":
			if err := AskCourse(c); err != nil {
				return err
			}
		case "faculty":
			user.Faculty = &value
			user.EducationForm = nil
			user.Course = nil
			user.StudyGroup = nil
			if err := AskForm(c); err != nil {
				return err
			}
		case "form":
			user.EducationForm = &value
			if err := AskCourse(c); err != nil {
				return err
			}
		case "course":
			user.Course = &value
			user.State = db.UserStateNone
			if err := AskGroup(c); err != nil {
				return err
			}
		case "group":
			groups, success, err := api.GetAllGroups()
			if !success {
				return err
			}
			LogError(c, err)

			user.StudyGroup = &value
			for _, group := range groups {
				if group.Group.Key == value {
					user.Faculty = &group.Faculty.Key
					user.EducationForm = &group.EducationForm.Key
					user.Course = &group.Course.Key
					break
				}
			}

			user.State = db.UserStateNone
			if err := c.Edit("Готово!"); err != nil {
				return err
			}
		case "discipline":
			if err := SendSubject(c, value); err != nil {
				return err
			}
			if err := DeleteMessage(c); err != nil {
				return err
			}
		case "update":
			if user.StudyGroup == nil {
				if err := AskGroup(c); err != nil {
					return err
				}
				return nil
			}
			date := time.Now()
			if value != "today" {
				var err error
				date, err = time.Parse("02.01.2006", value)
				if err != nil {
					return err
				}
			}

			if value2 != "" {
				msgIds := strings.Split(value2, ",")
				for _, id := range msgIds {
					msg := tele.StoredMessage{
						ChatID:    c.Chat().ID,
						MessageID: id,
					}

					if err := c.Bot().Delete(msg); err != nil {
						LogError(c, err)
					}
				}
			}

			if err := SendSchedule(c, c.Message(), true, FormatDay, date, date); err != nil {
				return err
			}
		case "settings":
			markup := &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{tele.InlineButton{Text: "Формат розкладу", Data: "settings/format"}},
					{tele.InlineButton{Text: "Фільтрація дисциплін", Data: "settings/disciplines"}},
					{tele.InlineButton{Text: "Закрити ❌", Data: "delete"}},
				},
			}
			if err := c.Edit("Налаштування", markup); err != nil {
				return err
			}
		case "settings/format":
			switch value {
			case "show-groups":
				user.Settings.ShowGroups = (user.Settings.ShowGroups + 1) % 3
			case "show-teacher":
				user.Settings.ShowTeacher = !user.Settings.ShowTeacher
			}

			showGroups := []string{"Не показувати", "Частково", "Всі"}[user.Settings.ShowGroups]
			showTeacher := []string{"Не показувати", "Показувати"}[boolToInt(user.Settings.ShowTeacher)]
			markup := &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{{Text: "Групи: " + showGroups, Data: "settings/format:show-groups"}},
					{{Text: "Викладач: " + showTeacher, Data: "settings/format:show-teacher"}},
					{{Text: "Назад ↩️", Data: "settings"}},
				},
			}
			if err := c.Edit("Формат повідомлення", markup); err != nil {
				return err
			}
		case "settings/disciplines":
			if user.StudyGroup == nil {
				if err := AskGroup(c); err != nil {
					return err
				}
				return nil
			}
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
			}
			if value == "all-on" {
				user.Settings.HiddenSubjects = []string{}
			}
			if value2 != "" {
				status := value2 == "on"
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
					}
				}
			}

			var keyboard [][]tele.InlineButton
			numSelected := len(subjects)

			for hashed, subject := range subjects {
				hidden := slices.Contains(user.Settings.HiddenSubjects, subject)
				prefix := []string{"❌", "✅"}[boolToInt(!hidden)]
				statusText := "off"
				if hidden {
					numSelected -= 1
					statusText = "on"
				}

				keyboard = append(keyboard, []tele.InlineButton{{
					Text: fmt.Sprintf("%s %s", prefix, subject),
					Data: "settings/disciplines:" + hashed + "|" + statusText,
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

			text := "Фільтрація дисциплін"
			if numSelected < 6 {
				alarm := ""
				if numSelected == 1 {
					alarm = fmt.Sprintf("⚠️ Увага! Ви вибрали тільки %d дисципліну", numSelected)
				} else if numSelected <= 4 && numSelected != 0 {
					alarm = fmt.Sprintf("⚠️ Увага! Ви вибрали тільки %d дисципліни", numSelected)
				} else {
					alarm = fmt.Sprintf("⚠️ Увага! Ви вибрали тільки %d дисциплін", numSelected)
				}
				text += "\n\n" + alarm
				keyboard = append(keyboard, []tele.InlineButton{{Text: alarm, Data: "-"}})
			}

			keyboard = append(keyboard, []tele.InlineButton{{Text: "Назад ↩️", Data: "settings"}})

			markup := &tele.ReplyMarkup{
				InlineKeyboard: keyboard,
			}
			if err := c.Edit(text, markup); err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
				return err
			}
		case "delete":
			if err := DeleteMessage(c); err != nil {
				return err
			}
		case "set-state":
			state, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			if err := DeleteMessage(c); err != nil {
				return err
			}
			user.State = state
		case "update-teacher":
			user.State = db.UserStateNone

			searches, err := user.GetSearches(db.UserSearchTypeEmployee)
			if err != nil {
				LogError(c, err)
			} else {
				searches = slices.DeleteFunc(searches, func(search string) bool {
					return search == value
				})
				searches = append(searches, value)
				if len(searches) > 6 {
					searches = searches[len(searches)-6:]
				}
				if err := user.SetSearches(db.UserSearchTypeEmployee, searches); err != nil {
					return err
				}
			}

			employee := api.KeyValue{
				Key:   value,
				Value: "",
			}

			date := time.Now()
			if value2 != "today" {
				var err error
				date, err = time.Parse("02.01.2006", value2)
				if err != nil {
					return err
				}
			}

			if value3 != "" {
				msgIds := strings.Split(value3, ",")
				for _, id := range msgIds {
					msg := tele.StoredMessage{
						ChatID:    c.Chat().ID,
						MessageID: id,
					}

					if err := c.Bot().Delete(msg); err != nil {
						LogError(c, err)
					}
				}
			}

			schedule, err := api.GetEmployeeSchedule(employee, date, date)
			if err != nil {
				return err
			}
			if len(schedule) == 0 {
				markup := GetTeacherDayMarkup(c, value, date.Format("02.01.2006"), nil)
				text := consts.WeekDays[date.Weekday()] + ", " + date.Format("02.01.2006")

				text += "\n\nРозклад порожній"
				if err := c.Edit(Escape(text), tele.ModeMarkdownV2, markup); err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
					return err
				}
			} else {
				day := &schedule[0]

				texts := FormatDay(c, day, false)
				messageIdsToDelete := []string{}
				for j, text := range texts {
					lastIteration := j == len(texts)-1

					var markup *tele.ReplyMarkup
					if lastIteration {
						markup = GetTeacherDayMarkup(c, value, day.Date, messageIdsToDelete)
					}

					msg, err := c.Bot().Edit(c.Message(), text, tele.ModeMarkdownV2, markup)
					if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
						return err
					}

					day.MessageIds = append(day.MessageIds, msg.ID)

					if !lastIteration {
						messageIdsToDelete = append(messageIdsToDelete, fmt.Sprintf("%d", msg.ID))
					}
				}
			}
		case "-":
			//
		default:
			LogError(c, fmt.Errorf("unknown query callback %v", c.Data()))
		}

		if err := user.Save(); err != nil {
			return err
		}
	}

	return nil
}
