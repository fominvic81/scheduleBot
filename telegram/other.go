package telegram

import tele "gopkg.in/telebot.v3"

func MainHandler(c tele.Context) error {
	keyboard := GetReplyKeyboard(KeyboardMain)

	if err := c.Delete(); err != nil {
		return err
	}

	_, err := c.Bot().Send(c.Recipient(), "Головна", &tele.ReplyMarkup{
		ReplyKeyboard:  keyboard,
		ResizeKeyboard: true,
	})

	return err
}

func OtherHandler(c tele.Context) error {
	keyboard := GetReplyKeyboard(KeyboardOther)

	if err := c.Delete(); err != nil {
		return err
	}

	_, err := c.Bot().Send(c.Recipient(), "Інше", &tele.ReplyMarkup{
		ReplyKeyboard:  keyboard,
		ResizeKeyboard: true,
	})

	return err
}
