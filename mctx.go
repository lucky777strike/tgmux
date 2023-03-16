package tgmux

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Ctx struct {
	Msg   *tgbotapi.Message
	Bot   *tgbotapi.BotAPI
	State *UserState
}

func (c *Ctx) SendErrorMessage(err error) {
	errorMessage := err.Error()
	reply := tgbotapi.NewMessage(c.Msg.Chat.ID, errorMessage)
	reply.ReplyToMessageID = c.Msg.MessageID
	_, sendErr := c.Bot.Send(reply)
	if sendErr != nil {
		log.Printf("Error sending error message: %v\n", sendErr)
	}
}
