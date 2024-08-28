package telegram

import tele "gopkg.in/telebot.v3"

type KeyboardButton struct {
	Text    string
	Handler func(tele.Context) error
}

type Keyboard [][]KeyboardButton

func GetKeyboard() Keyboard {
	return Keyboard{
		{
			{Text: "Cьогодні", Handler: Day},
			{Text: "Завтра", Handler: Next},
			{Text: "Післязавтра", Handler: NextNext},
		},
		{
			{Text: "Тиждень", Handler: Week},
			{Text: "Наступний тиждень", Handler: NextWeek},
		},
		{
			{Text: "Стисло(два тижні)", Handler: Short},
			{Text: "Знайти предмет", Handler: Subject},
		},
	}
}

func GetKeyboardMarkup() *tele.ReplyMarkup {
	keyboard := GetKeyboard()
	replyKeyboard := make([][]tele.ReplyButton, len(keyboard))

	for i, row := range keyboard {
		replyKeyboard[i] = make([]tele.ReplyButton, len(row))
		for j, button := range row {
			replyKeyboard[i][j] = tele.ReplyButton{
				Text: button.Text,
			}
		}
	}

	return &tele.ReplyMarkup{
		ReplyKeyboard:  replyKeyboard,
		ResizeKeyboard: true,
	}
}
