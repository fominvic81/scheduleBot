package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func Start(c tele.Context) error {
	err := c.Send(HelpMsg(), c.Get("keyboard").(func() *tele.ReplyMarkup)())

	if err != nil {
		return err
	}

	_, err = Ask(c)
	if err != nil {
		return err
	}

	return nil
}
