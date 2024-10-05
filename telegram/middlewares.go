package telegram

import (
	"database/sql"
	"errors"

	"github.com/fominvic81/scheduleBot/consts"
	"github.com/fominvic81/scheduleBot/db"

	tele "gopkg.in/telebot.v3"
)

func EmptyStrAsNil(str string) *string {
	if len(str) == 0 {
		return nil
	}
	return &str
}

func UserMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		database := c.Get("database").(*sql.DB)

		sender := c.Sender()
		if sender == nil {
			return errors.New("failed to identify user")
		}
		user, err := db.GetOrCreateUser(database, sender.ID, sender.FirstName)
		if err != nil {
			return err
		}
		user.Firstname = sender.FirstName
		user.Lastname = EmptyStrAsNil(sender.LastName)
		user.Username = EmptyStrAsNil(sender.Username)
		user.Messages += 1
		err = user.Save()

		if err != nil {
			return err
		}

		c.Set("user", user)
		return next(c)
	}
}

func LogMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		LogAction(c)
		return next(c)
	}
}

func KeyboardMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user := c.Get("user").(*db.User)

		if user.KeyboardVersion != consts.KeyboardVersion {
			keyboard := GetReplyKeyboard()

			_, err := c.Bot().Send(c.Recipient(), "Оновлено клавіатуру", &tele.ReplyMarkup{
				ReplyKeyboard:  keyboard,
				ResizeKeyboard: true,
			})
			if err != nil {
				return err
			}

			user.KeyboardVersion = consts.KeyboardVersion

			err = user.Save()
			if err != nil {
				return err
			}
		}

		return next(c)
	}
}

func MetricWriteMessage(metric *db.Metric, message *tele.Message) {
	metric.Content = message.Text

	media := message.Media()
	if media != nil {
		mediaType := media.MediaType()
		metric.MediaType = mediaType
		mediaFile := media.MediaFile()

		if mediaFile != nil {
			metric.MediaId = mediaFile.FileID
		}

		if metric.Content != "" { // is it even possible?
			metric.Content += " |> " + message.Caption
		} else {
			metric.Content = message.Caption
		}
	}
	metric.AlbumId = message.AlbumID

	if message.ReplyTo != nil {
		metric.ReplyTo = int64(message.ReplyTo.ID)
	}
}

func MetricsMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user := c.Get("user").(*db.User)

		metric := db.Metric{}

		sender := c.Sender()
		chat := c.Chat()
		callback := c.Callback()
		message := c.Update().Message
		editedMessage := c.Update().EditedMessage

		if sender != nil {
			metric.UserId = sender.ID
		}
		if chat != nil {
			metric.ChatId = chat.ID
		}

		if callback != nil {
			metric.EventType = db.EventTypeReplyCallback

			metric.Content = callback.Data
			if callback.Message != nil {
				metric.ReplyTo = int64(callback.Message.ID)
			}
		} else if editedMessage != nil {
			metric.EventType = db.EventTypeMessageEdited

			MetricWriteMessage(&metric, editedMessage)
		} else if message != nil {
			metric.EventType = db.EventTypeMessage

			MetricWriteMessage(&metric, message)
		}

		if user.Messages == 1 {
			metric.Flags |= db.MetricFlagFirstMessage
		}

		database := c.Get("database").(*sql.DB)

		err := db.WriteMetric(database, metric)
		if err != nil {
			LogError(err, c)
		}

		return next(c)
	}
}
