package telegram

import (
	"database/sql"
	"log"

	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func ErrorHanler(err error, c tele.Context) {
	log.Println(err)
	err2 := c.Send("Сталася помилка")
	if err2 != nil {
		log.Println(err2)
	}

	admins, err2 := db.GetAdminUsers(c.Get("database").(*sql.DB))
	if err2 != nil {
		log.Println(err2)
	} else {
		for _, admin := range admins {
			_, _ = c.Bot().Send(tele.ChatID(admin.Id), "Error: "+err.Error())
		}
	}
}
