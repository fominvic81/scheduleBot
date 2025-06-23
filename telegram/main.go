package telegram

import (
	"database/sql"
	"log"
	"time"

	"github.com/fominvic81/scheduleBot/api"
	tele "gopkg.in/telebot.v3"
)

func preload() {
	if _, _, err := api.GetAllEmployeesAndChairs(); err != nil {
		log.Println(err)
	}
	if _, _, err := api.GetAllStudyGroups(); err != nil {
		log.Println(err)
	}
}

func setCommands(b *tele.Bot) {
	commands := GetCommands()

	teleCommands := make([]tele.Command, 0, len(commands))
	for _, command := range commands {
		teleCommands = append(teleCommands, tele.Command{
			Text:        command.Text,
			Description: command.Description,
		})
	}
	err := b.SetCommands(teleCommands)
	if err != nil {
		log.Println(err)
	}
}

func Init(token string, database *sql.DB) {
	go preload()

	pref := tele.Settings{
		Token:   token,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: ErrorHandler,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	setCommands(b)

	b.Use(RecoverMiddleware)
	b.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set("database", database)
			return next(c)
		}
	})
	b.Use(LogMiddleware)
	b.Use(UserMiddleware)
	b.Use(MetricsMiddleware)
	// b.Use(BanMiddleware)
	b.Use(KeyboardMiddleware)

	commands := GetCommands()

	for _, command := range commands {
		b.Handle("/"+command.Text, command.Handler)
	}
	b.Handle("/send", SendHandler)

	b.Handle(tele.OnText, TextHandler)
	b.Handle(tele.OnCallback, CallbackDataHandler)

	// Fix middlewares not being called if there are no matching handlers
	b.Handle(tele.OnMedia, func(ctx tele.Context) error {
		return nil
	})
	b.Handle(tele.OnEdited, func(ctx tele.Context) error {
		return nil
	})

	b.Start()
}
