package telegram

import (
	"fmt"

	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func GroupInRows(buttons []tele.Btn, perRow int) []tele.Row {
	numRows := (len(buttons) + perRow - 1) / perRow
	rows := make([]tele.Row, 0, numRows+1)
	for i := range numRows {
		rows = append(rows, buttons[i*perRow:min((i+1)*perRow, len(buttons))])
	}
	return rows
}

func AskFaculty(c tele.Context) error {
	user := c.Get("user").(*db.User)

	filters, err := api.GetFilters()
	if err != nil {
		return err
	}

	user.State = db.UserStateSearchGroup
	if err := user.Save(); err != nil {
		return err
	}

	buttons := make([]tele.Btn, 0, len(filters.Faculties))
	for _, entry := range filters.Faculties {
		buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "faculty:"+entry.Key))
	}

	rows := GroupInRows(buttons, 3)
	rows = append(rows, tele.Row{tele.Btn{
		Text: "Скасувати ❌",
		Data: fmt.Sprintf("set-state:%v", db.UserStateNone),
	}})
	markup := tele.ReplyMarkup{}
	markup.Inline(rows...)

	if err = c.EditOrSend("Виберіть факультет або введіть шифр групи(наприклад ІПЗ-32 або ІПЗ)", &markup); err != nil {
		return err
	}

	return nil
}

func AskForm(c tele.Context) error {
	user := c.Get("user").(*db.User)

	filters, err := api.GetFilters()
	if err != nil {
		return err
	}

	user.State = db.UserStateSearchGroup
	if err := user.Save(); err != nil {
		return err
	}

	buttons := make([]tele.Btn, 0, len(filters.EducationForms))
	for _, entry := range filters.EducationForms {
		buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "form:"+entry.Key))
	}

	rows := GroupInRows(buttons, 3)
	rows = append(rows, tele.Row{tele.Btn{
		Text: "Назад ↩️",
		Data: "ask-faculty",
	}})
	markup := tele.ReplyMarkup{}
	markup.Inline(rows...)

	if err = c.EditOrSend("Виберіть форму навчання або введіть шифр групи(наприклад ІПЗ-32 або ІПЗ)", &markup); err != nil {
		return err
	}

	return nil
}

func AskCourse(c tele.Context) error {
	user := c.Get("user").(*db.User)

	filters, err := api.GetFilters()
	if err != nil {
		return err
	}

	user.State = db.UserStateSearchGroup
	if err := user.Save(); err != nil {
		return err
	}

	buttons := make([]tele.Btn, 0, len(filters.Courses))
	for _, entry := range filters.Courses {
		buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "course:"+entry.Key))
	}

	rows := GroupInRows(buttons, 4)
	rows = append(rows, tele.Row{tele.Btn{
		Text: "Назад ↩️",
		Data: "ask-form",
	}})
	markup := tele.ReplyMarkup{}
	markup.Inline(rows...)

	if err = c.EditOrSend("Виберіть курс або введіть шифр групи(наприклад ІПЗ-32 або ІПЗ)", &markup); err != nil {
		return err
	}

	return nil
}

func AskGroup(c tele.Context) error {
	user := c.Get("user").(*db.User)

	if user.Faculty == nil {
		return AskFaculty(c)
	}
	if user.EducationForm == nil {
		return AskForm(c)
	}
	if user.Course == nil {
		return AskCourse(c)
	}

	studyGroups, err := api.GetStudyGroups(*user.Faculty, *user.EducationForm, *user.Course)
	if err != nil {
		return err
	}

	if len(studyGroups) == 0 {
		if err = c.EditOrSend("Не знайдено навчальних груп з такими даними", &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{{{
				Text: "Назад ↩️",
				Data: "ask-course",
			}}},
		}); err != nil {
			return err
		}
		return nil
	}

	user.State = db.UserStateSearchGroup
	if err = user.Save(); err != nil {
		return err
	}

	buttons := make([]tele.Btn, 0, len(studyGroups))
	for _, entry := range studyGroups {
		buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "group:"+entry.Key))
	}

	rows := GroupInRows(buttons, 3)
	rows = append(rows, tele.Row{tele.Btn{
		Text: "Назад ↩️",
		Data: "ask-course",
	}})
	markup := tele.ReplyMarkup{}
	markup.Inline(rows...)

	if err = c.EditOrSend("Виберіть навчальну групу або введіть шифр групи(наприклад ІПЗ-32 або ІПЗ)", &markup); err != nil {
		return err
	}

	return err
}
