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
	botToken := "TOKEN"

	handler, err := tgmux.NewHandler(botToken)
	if err != nil {
		log.Panic(err)
	}
	handler.SetCustomMessages(&tgmux.Messages{NoCommand: "Данная команда недоступна",
		InternalError: "Внутренняя ошибка сервера"})

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
