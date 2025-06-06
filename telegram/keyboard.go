package telegram

import tele "gopkg.in/telebot.v3"

type KeyboardButton struct {
	Text    string
	Handler func(tele.Context) error
}

type Keyboard [][]KeyboardButton

const (
	KeyboardMain = iota
	KeyboardOther
)

func GetKeyboards() []Keyboard {
	return []Keyboard{
		{
			{
				{Text: "Cьогодні", Handler: DayHandler},
				{Text: "Завтра", Handler: NextHandler},
				{Text: "Післязавтра", Handler: NextNextHandler},
			},
			{
				{Text: "Тиждень", Handler: WeekHandler},
				{Text: "Наступний тиждень", Handler: NextWeekHandler},
			},
			{
				{Text: "Інше", Handler: OtherHandler},
			},
		},
		{
			{
				{Text: "Стисло(два тижні)", Handler: ShortHandler},
				{Text: "Знайти предмет", Handler: SubjectHandler},
			},
			{
				{Text: "Змінити групу", Handler: SetGroupHandler},
				{Text: "Змінити факультет і групу", Handler: SetDataHandler},
			},
			{
				{Text: "Налаштування", Handler: SettingsHandler},
			},
			{
				{Text: "Головна", Handler: MainHandler},
			},
		},
	}
}

func GetReplyKeyboard(keyboardIndex int) [][]tele.ReplyButton {
	keyboard := GetKeyboards()[keyboardIndex]
	replyKeyboard := make([][]tele.ReplyButton, len(keyboard))

	for i, row := range keyboard {
		replyKeyboard[i] = make([]tele.ReplyButton, len(row))
		for j, button := range row {
			replyKeyboard[i][j] = tele.ReplyButton{
				Text: button.Text,
			}
		}
	}

	return replyKeyboard
}
