package tgmux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Ctx struct {
	Msg   *tgbotapi.Message
	Bot   *tgbotapi.BotAPI
	State UserStateInterface
	Log   Logger
}

func (c *Ctx) SendErrorMessage(err error) {
	errorMessage := err.Error()
	reply := tgbotapi.NewMessage(c.Msg.Chat.ID, errorMessage)
	reply.ReplyToMessageID = c.Msg.MessageID
	_, sendErr := c.Bot.Send(reply)
	if sendErr != nil {
		c.Log.Println(err)
	}
}
