package telegram

import (
	"database/sql"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func SetCommands(b *tele.Bot) {
	commands := GetCommands()

	teleCommands := make([]tele.Command, 0, len(commands))
	for _, command := range commands {
		teleCommands = append(teleCommands, tele.Command{
			Text:        command.Text,
			Description: command.Description,
		})
	}
	b.SetCommands(teleCommands)
}

func Init(token string, database *sql.DB) {
	pref := tele.Settings{
		Token:   token,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: ErrorHanler,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	SetCommands(b)

	b.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set("database", database)
			return next(c)
		}
	})
	b.Use(UserMiddleware)
	b.Use(KeyboardMiddleware)

	commands := GetCommands()

	for _, command := range commands {
		b.Handle("/"+command.Text, command.Handler)
	}

	b.Handle(tele.OnText, HandleText)
	b.Handle(tele.OnCallback, CallbackData)

	b.Start()
}
