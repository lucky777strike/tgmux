package tgmux

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Test() string {
	return "Test passed"
}

type TgHandler struct {
	bot        *tgbotapi.BotAPI
	croutes    map[string]func(*Ctx) //commandroutes
	sroutes    map[string]func(*Ctx)
	userStates *UserStateManager
	log        Logger
	messages   *Messages
}

func NewHandler(token string) (*TgHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TgHandler{bot: bot,
			croutes:    make(map[string]func(*Ctx)),
			sroutes:    make(map[string]func(*Ctx)),
			userStates: NewUserStateManager(),
			log:        log.Default(),
			messages:   defaultMessages},

		nil
}

func NewHandlerWithLogger(token string, logger Logger) (*TgHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TgHandler{bot: bot,
			croutes:    make(map[string]func(*Ctx)),
			sroutes:    make(map[string]func(*Ctx)),
			userStates: NewUserStateManager(),
			log:        logger},
		nil
}

func (t *TgHandler) HandleCmd(command string, f func(*Ctx)) {
	t.croutes[command] = f
}
func (t *TgHandler) HandleState(command string, f func(*Ctx)) {
	t.sroutes[command] = f
}

func (t *TgHandler) processUpdate(update *tgbotapi.Update) {
	if update.Message != nil {
		userID := update.Message.From.ID
		userState := t.userStates.GetUserState(int64(userID))
		mctx := &Ctx{update.Message, t.bot, userState, t.log}

		currentFunction := userState.GetCurrentFunction()
		if currentFunction != "" {
			handler, ok := t.sroutes[currentFunction]
			if ok {
				go handler(mctx)
			} else {
				t.userStates.ResetUserFunction(int64(userID))
				t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, t.messages.InternalError))
				errorMsg := fmt.Sprintf("State handler not found for user %d, state function: %s", userID, currentFunction)
				t.log.Println(errors.New(errorMsg))
			}
			return
		}
		handler, ok := t.croutes[update.Message.Text]
		if ok {
			go handler(mctx)
		} else {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, t.messages.NoCommand))
		}
	}
}

func (t *TgHandler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		t.processUpdate(&update)
	}
}

func (t *TgHandler) SetCustomMessages(messages *Messages) {
	t.messages = messages
}
