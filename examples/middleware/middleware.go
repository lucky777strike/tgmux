package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lucky777strike/tgmux"
)

func main() {
	botToken := "5990324330:AAEZdIaNzVTSQIlZJnU9zwj1QhfnPSDXr5g"

	handler, err := tgmux.NewHandler(botToken)
	if err != nil {
		log.Panic(err)
	}
	handler.AddMiddleware(sendInlineMenu)
	handler.HandleCmd("/start", startCommand)
	handler.HandleCmd("/sum", sumCommand)
	handler.HandleState("sum", sumCommand)
	handler.Start()
}

func startCommand(c *tgmux.Ctx) {
	welcomeMessage := fmt.Sprintf("Hello, %s! Welcome to the example bot.", c.Msg.From.FirstName)

	reply := tgbotapi.NewMessage(c.Msg.Chat.ID, welcomeMessage)
	reply.ReplyToMessageID = c.Msg.MessageID

	_, err := c.Bot.Send(reply)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func sumCommand(c *tgmux.Ctx) {
	currentFunction := c.State.GetCurrentFunction()
	if currentFunction == "" {
		c.State.SetCurrentFunction("sum")
		reply := tgbotapi.NewMessage(c.Msg.Chat.ID, "Please send number one")
		reply.ReplyToMessageID = c.Msg.MessageID
		_, err := c.Bot.Send(reply)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	} else if _, ok := c.State.GetData()["first"]; !ok {
		number1, err := strconv.Atoi(c.Msg.Text)
		if err != nil {
			c.SendErrorMessage(errors.New("Invalid input. Please send a valid integer."))
			return
		}
		c.State.UpdateData("first", number1)
		reply := tgbotapi.NewMessage(c.Msg.Chat.ID, "Please send number two")
		reply.ReplyToMessageID = c.Msg.MessageID
		_, err = c.Bot.Send(reply)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	} else if _, ok := c.State.GetData()["second"]; !ok {
		number2, err := strconv.Atoi(c.Msg.Text)
		if err != nil {
			c.SendErrorMessage(errors.New("Invalid input. Please send a valid integer."))
			return
		}
		c.State.UpdateData("second", number2)

		data := c.State.GetData()
		number1 := data["first"].(int)
		sum := number1 + number2
		reply := tgbotapi.NewMessage(c.Msg.Chat.ID, fmt.Sprintf("The sum of the two numbers is: %d", sum))
		reply.ReplyToMessageID = c.Msg.MessageID
		_, err = c.Bot.Send(reply)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		c.State.SetCurrentFunction("")
		c.State.SetData(make(map[string]interface{}))
	}
}
func sendInlineMenu(ctx *tgmux.Ctx) {
	menuText := "Please choose an option:"

	inlineMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Start", "/start"),
			tgbotapi.NewInlineKeyboardButtonData("Help", "/help"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Some Command", "/some_command"),
			tgbotapi.NewInlineKeyboardButtonData("Reset", "/reset"),
		),
	)

	msg := tgbotapi.NewMessage(ctx.Msg.Chat.ID, menuText)
	msg.ReplyMarkup = inlineMenu
	ctx.Bot.Send(msg)
}
