package telegram

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func SearchGroupHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)

	groups, success, err := api.GetAllGroups()
	if !success {
		return err
	}
	LogError(c, err)

	results := []api.Group{}
	exactResults := []api.Group{}

	query := strings.ReplaceAll(strings.ToLower(c.Text()), " ", "-")

	for _, group := range groups {
		if strings.ToLower(group.Group.Value) == query {
			exactResults = append(exactResults, group)
		}
		if strings.Contains(strings.ToLower(group.Group.Value), query) {
			results = append(results, group)
		}
		if len(results) > 50 {
			break
		}
	}

	if len(results) == 0 {
		markup := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{{{
				Text: "Скасувати пошук ❌",
				Data: fmt.Sprintf("set-state:%v", db.UserStateNone),
			}}},
		}
		if err := c.Send("Такої групи не знайдено", markup); err != nil {
			return err
		}
		return nil
	}

	if len(exactResults) == 1 && len(results) == 1 {
		group := exactResults[0]
		user.State = db.UserStateNone

		user.Faculty = &group.Faculty.Key
		user.EducationForm = &group.EducationForm.Key
		user.Course = &group.Course.Key
		user.StudyGroup = &group.Group.Key

		if err = user.Save(); err != nil {
			return err
		}
		err = c.Send("Встановлено групу " + group.Group.Value + " (" + group.Faculty.Value + ", " + group.EducationForm.Value + ", " + group.Course.Key + " курс)")
		if err != nil {
			return err
		}
		return nil
	}

	buttons := [][]tele.InlineButton{}
	for _, group := range results {
		buttons = append(buttons, []tele.InlineButton{{
			Text: group.Group.Value + " (" + group.Faculty.Value + ", " + group.EducationForm.Value + ", " + group.Course.Key + " курс)",
			Data: "group:" + group.Group.Key,
		}})
	}

	markup := &tele.ReplyMarkup{
		InlineKeyboard: buttons,
	}
	if err = c.Send("Виберіть групу ", markup); err != nil {
		return err
	}

	return nil
}

func SetSearchTeacherHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)

	user.State = db.UserStateSearchTeacher
	if err := user.Save(); err != nil {
		return nil
	}

	markup := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{{{
			Text: "Скасувати пошук ❌",
			Data: fmt.Sprintf("set-state:%v", db.UserStateNone),
		}}},
	}
	if err := c.Send("Введіть ім'я викладача", markup); err != nil {
		return err
	}

	return nil
}

func SearchTeacherHandler(c tele.Context) error {
	employeesAndChairs, success, err := api.GetAllEmployeesAndChairs()
	if !success {
		return err
	}
	LogError(c, err)

	results := []api.KeyValue{}

	query := strings.ToLower(c.Text())
	words := strings.Split(strings.ReplaceAll(query, ".", " "), " ")
	words = slices.DeleteFunc(words, func(word string) bool {
		return word == ""
	})

	for _, employee := range employeesAndChairs.Employees {
		name := strings.ToLower(employee.Value)
		nameWords := strings.Split(name, " ")
		if strings.Contains(name, query) {
			results = append(results, employee)
			continue
		}

		if len(words) <= len(nameWords) {
			matches := true
			for i, word := range words {
				if !strings.HasPrefix(nameWords[i], word) {
					matches = false
					break
				}
			}
			if matches {
				results = append(results, employee)
				continue
			}
		}
	}

	if len(results) == 0 {
		markup := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{{{
				Text: "Скасувати пошук ❌",
				Data: fmt.Sprintf("set-state:%v", db.UserStateNone),
			}}},
		}
		if err := c.Send("Жодного викладача не знайдено", markup); err != nil {
			return err
		}
		return nil
	}

	buttons := [][]tele.InlineButton{}
	for _, employee := range results {
		buttons = append(buttons, []tele.InlineButton{{
			Text: employee.Value,
			Data: "update-teacher:" + employee.Key + "|" + time.Now().Format("02.01.2006"),
		}})
		if len(buttons) >= 50 {
			buttons = append(buttons, []tele.InlineButton{{
				Text: "...",
				Data: "-",
			}})
			break
		}
	}
	buttons = append(buttons, []tele.InlineButton{{
		Text: "Скасувати пошук ❌",
		Data: fmt.Sprintf("set-state:%v", db.UserStateNone),
	}})

	markup := &tele.ReplyMarkup{
		InlineKeyboard: buttons,
	}

	if err = c.Send("Виберіть викладача ", markup); err != nil {
		return err
	}

	return nil
}
