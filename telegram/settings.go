package telegram

import tele "gopkg.in/telebot.v3"

func SettingsHandler(c tele.Context) error {
	return c.Send("Налаштування", &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{tele.InlineButton{Text: "Формат розкладу", Data: "settings/format"}},
			{tele.InlineButton{Text: "Фільтрація дисциплін", Data: "settings/disciplines"}},
			{tele.InlineButton{Text: "Закрити", Data: "delete"}},
		},
	})
}
