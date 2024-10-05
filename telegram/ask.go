package telegram

import (
	"github.com/fominvic81/scheduleBot/api"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func GroupInRows(buttons []tele.Btn, perRow int) tele.ReplyMarkup {
	markup := tele.ReplyMarkup{}
	numRows := (len(buttons) + perRow - 1) / perRow
	rows := make([]tele.Row, 0, numRows)
	for i := range numRows {
		rows = append(rows, buttons[i*perRow:min((i+1)*perRow, len(buttons))])
	}
	markup.Inline(rows...)
	return markup
}

func Ask(c tele.Context) (bool, error) {
	user := c.Get("user").(*db.User)

	if user.Faculty == nil || user.EducationForm == nil || user.Course == nil {
		filters, err := api.GetFilters()
		if err != nil {
			return false, nil
		}

		err = nil
		if user.Faculty == nil {
			buttons := make([]tele.Btn, 0, len(filters.Faculties))
			for _, entry := range filters.Faculties {
				buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "faculty:"+entry.Key))
			}

			markup := GroupInRows(buttons, 3)
			err = c.Send("Виберіть факультет", &markup)
		} else if user.EducationForm == nil {
			buttons := make([]tele.Btn, 0, len(filters.EducationForms))
			for _, entry := range filters.EducationForms {
				buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "form:"+entry.Key))
			}

			markup := GroupInRows(buttons, 3)
			err = c.Send("Виберіть форму навчання", &markup)
		} else if user.Course == nil {
			buttons := make([]tele.Btn, 0, len(filters.Courses))
			for _, entry := range filters.Courses {
				buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "course:"+entry.Key))
			}

			markup := GroupInRows(buttons, 3)
			err = c.Send("Виберіть курс", &markup)
		}

		return true, err
	} else if user.StudyGroup == nil {
		studyGroups, err := api.GetStudyGroups(*user.Faculty, *user.EducationForm, *user.Course)

		if err != nil {
			return false, err
		}

		if len(studyGroups) == 0 {
			user.Faculty = nil
			user.EducationForm = nil
			user.Course = nil
			err = user.Save()
			if err != nil {
				return false, err
			}
			return false, c.Send("Не знайдено навчальних груп з такими даними")
		}

		buttons := make([]tele.Btn, 0, len(studyGroups))
		for _, entry := range studyGroups {
			buttons = append(buttons, (&tele.ReplyMarkup{}).Data(entry.Value, "group:"+entry.Key))
		}

		markup := GroupInRows(buttons, 3)

		return true, c.Send("Виберіть навчальну групу", &markup)
	}

	return false, nil
}
