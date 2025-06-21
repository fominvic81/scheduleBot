package telegram

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func LogAction(c tele.Context) {
	text := ""

	sender := c.Sender()
	if sender != nil {
		text += fmt.Sprintf("[%d @%s] ", sender.ID, sender.Username)
	}

	callback := c.Callback()
	if callback != nil {
		text += fmt.Sprintf("cb{%s: %s} ", callback.MessageID, callback.Data)
	}

	message := c.Update().Message
	if message != nil {
		if message.Text != "" {
			text += fmt.Sprintf("msg{%d: %s} ", message.ID, message.Text)
		}
		if message.Media() != nil {
			text += "(Media) "
		}
	}

	log.Println(text)
}

func LogError(err error, c tele.Context) {
	if err == nil {
		return
	}
	errorText := ""

	sender := c.Sender()
	if sender != nil {
		errorText = fmt.Sprintf("[%d @%s] Error: %s", sender.ID, sender.Username, err.Error())
	} else {
		errorText = fmt.Sprintf("[] Error: %s", err.Error())
	}

	log.Println(errorText)

	admins, err2 := db.GetAdminUsers(c.Get("database").(*sql.DB))
	if err2 != nil {
		log.Println(err2)
	} else {
		for _, admin := range admins {
			_, _ = c.Bot().Send(tele.ChatID(admin.Id), errorText)
		}
	}
}
