package telegram

import (
	"database/sql"
	"errors"

	"github.com/fominvic81/scheduleBot/db"
	tele "gopkg.in/telebot.v3"
)

func TextHandler(c tele.Context) error {
	user := c.Get("user").(*db.User)
	database := c.Get("database").(*sql.DB)

	keyboards := GetKeyboards()
	for _, keyboard := range keyboards {
		for _, row := range keyboard {
			for _, btn := range row {
				if c.Text() == btn.Text {
					if user.State != db.UserStateNone {
						user.State = db.UserStateNone
						if err := user.Save(); err != nil {
							return err
						}
					}
					if err := btn.Handler(c); err != nil {
						return err
					}
					return nil
				}
			}
		}
	}

	switch user.State {
	case db.UserStateSearchGroup:
		if err := SearchGroupHandler(c); err != nil {
			return err
		}
		return nil
	case db.UserStateSearchTeacher:
		if err := SearchTeacherHandler(c); err != nil {
			return err
		}
		return nil
	}

	if user.IsAdmin {
		replyTo := c.Message().ReplyTo

		var originalChatId *int64 = nil
		var originalMessageId *int = nil

		if replyTo != nil {
			result := database.QueryRow("SELECT original_chat_id, original_message_id FROM forwarded_messages WHERE chat_id = ? AND message_id = ?", replyTo.Chat.ID, replyTo.ID)
			err := result.Scan(&originalChatId, &originalMessageId)
			if errors.Is(err, sql.ErrNoRows) {
				if replyTo.Origin != nil {
					originalChatId = &replyTo.Origin.Chat.ID
					originalMessageId = &replyTo.Origin.MessageID
				}
			} else if err != nil {
				return err
			}

		}
		if originalChatId != nil && originalMessageId != nil {
			if _, err := c.Bot().Copy(tele.ChatID(*originalChatId), c.Message(), &tele.ReplyParams{MessageID: *originalMessageId}); err != nil {
				return err
			}
		}
	} else {
		admins, err := db.GetAdminUsers(c.Get("database").(*sql.DB))
		if err != nil {
			return err
		}
		for _, admin := range admins {
			message, err := c.Bot().Forward(tele.ChatID(admin.Id), c.Message())
			if err != nil {
				return err
			}

			query := "INSERT INTO forwarded_messages (chat_id, message_id, original_chat_id, original_message_id) VALUES (?, ?, ?, ?)"
			if _, err := database.Exec(query, admin.Id, message.ID, c.Chat().ID, c.Message().ID); err != nil {
				return err
			}
		}
	}

	return nil
}
