package telegram

import tele "gopkg.in/telebot.v3"

func HandleText(c tele.Context) error {

	keyboards := GetKeyboards()
	for _, keyboard := range keyboards {
		for _, row := range keyboard {
			for _, btn := range row {
				if c.Text() == btn.Text {
					return btn.Handler(c)
				}
			}
		}
	}

	return nil
}
