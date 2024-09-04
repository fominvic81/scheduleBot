package telegram

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func LogError(err error, c tele.Context) {
	sender := c.Sender()

	errorText := ""

	if sender != nil {
		errorText = fmt.Sprintf("[%d @%s] %s", sender.ID, sender.Username, err.Error())
	} else {
		errorText = err.Error()
	}

	log.Println(errorText)

	admins, err2 := db.GetAdminUsers(c.Get("database").(*sql.DB))
	if err2 != nil {
		log.Println(err2)
	} else {
		for _, admin := range admins {
			_, _ = c.Bot().Send(tele.ChatID(admin.Id), "Error: "+errorText)
		}
	}
}

func ErrorHanler(err error, c tele.Context) {
	LogError(err, c)

	err2 := c.Send("Сталася помилка")
	if err2 != nil {
		log.Println(err2)
	}
}
