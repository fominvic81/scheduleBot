package telegram

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type Command struct {
	Text        string
	Description string
	Handler     func(tele.Context) error
}

func GetCommands() []Command {
	return []Command{
		{Text: "start", Description: "Старт", Handler: Start},

		{Text: "day", Description: "Розклад на сьогодні", Handler: Day},
		{Text: "next", Description: "Розклад на завтра", Handler: Next},
		{Text: "nextnext", Description: "Розклад на післязавтра", Handler: NextNext},

		{Text: "week", Description: "Розклад на тиждень", Handler: Week},
		{Text: "nextweek", Description: "Розклад на наступний тиждень", Handler: NextWeek},

		{Text: "short", Description: "Стисло(два тижні)", Handler: Short},
		{Text: "subject", Description: "Знайти предмет", Handler: Subject},

		{Text: "setgroup", Description: "Змінити групу", Handler: SetGroup},
		{Text: "setdata", Description: "Змінити дані", Handler: SetData},
	}
}

func HelpMsg() string {
	msg := "Я бот, що дозволяє зручно та швидко слідкувати за розкладом занять в ЛНТУ\n" +
		"Основні команди:\n"

	for _, command := range GetCommands() {
		msg += fmt.Sprintf("/%s - %s\n", command.Text, command.Description)
	}

	return msg
}
