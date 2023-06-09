package tgmux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Ctx structure displays the context of the dialog (IT IS NOT context.Context).
type Ctx struct {
	Msg   *tgbotapi.Message
	Bot   *tgbotapi.BotAPI
	State UserStateInterface
	Log   Logger
}

// SendErrorMessage sends an error message as a reply to the user's message.
// If sending the error message fails, the error is logged.
func (c *Ctx) SendErrorMessage(err error) {
	errorMessage := err.Error()
	reply := tgbotapi.NewMessage(c.Msg.Chat.ID, errorMessage)
	reply.ReplyToMessageID = c.Msg.MessageID
	_, sendErr := c.Bot.Send(reply)
	if sendErr != nil {
		c.Log.Println(err)
	}
}
