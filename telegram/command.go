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
		{Text: "start", Description: "Старт", Handler: StartHandler},

		{Text: "day", Description: "Розклад на сьогодні", Handler: DayHandler},
		{Text: "next", Description: "Розклад на завтра", Handler: NextHandler},
		{Text: "nextnext", Description: "Розклад на післязавтра", Handler: NextNextHandler},

		{Text: "week", Description: "Розклад на тиждень", Handler: WeekHandler},
		{Text: "nextweek", Description: "Розклад на наступний тиждень", Handler: NextWeekHandler},

		{Text: "short", Description: "Стисло(два тижні)", Handler: ShortHandler},
		{Text: "subject", Description: "Знайти предмет", Handler: SubjectHandler},

		{Text: "setgroup", Description: "Змінити групу", Handler: AskGroup},
		{Text: "setdata", Description: "Змінити факультет і групу", Handler: AskFaculty},
		{Text: "teacher", Description: "Знайти викладача", Handler: SetSearchTeacherHandler},
		{Text: "mygroup", Description: "Моя група", Handler: MyGroupHandler},

		{Text: "settings", Description: "Налаштування", Handler: SettingsHandler},
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
