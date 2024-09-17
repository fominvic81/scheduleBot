package telegram

import (
	"log"

	tele "gopkg.in/telebot.v3"
)

func ErrorHandler(err error, c tele.Context) {
	LogError(err, c)

	err2 := c.Send("Сталася помилка")
	if err2 != nil {
		log.Println(err2)
	}
}
