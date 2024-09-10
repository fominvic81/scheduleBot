package telegram

import tele "gopkg.in/telebot.v3"

func Settings(c tele.Context) error {
	return c.Send("Налаштування", GetMarkup(c, &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{tele.InlineButton{Text: "Формат розкладу", Data: "settings/format"}},
			{tele.InlineButton{Text: "Вибіркові дисципліни", Data: "settings/disciplines"}},
			{tele.InlineButton{Text: "Закрити", Data: "delete"}},
		},
	}))
}
