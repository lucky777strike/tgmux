package tgmux

import (
	"context"
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgHandler struct {
	bot        *tgbotapi.BotAPI
	croutes    map[string]func(*Ctx) //command routes
	sroutes    map[string]func(*Ctx) //state routes
	userStates UserStateManagerInterface
	log        Logger
	messages   *Messages
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewHandler(token string) (*TgHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &TgHandler{bot: bot,
			croutes:    make(map[string]func(*Ctx)),
			sroutes:    make(map[string]func(*Ctx)),
			userStates: NewUserStateManager(),
			log:        log.Default(),
			messages:   defaultMessages,
			ctx:        ctx,
			cancel:     cancel},

		nil
}

func NewHandlerWithContext(ctx context.Context, cancel context.CancelFunc, token string) (*TgHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TgHandler{bot: bot,
			croutes:    make(map[string]func(*Ctx)),
			sroutes:    make(map[string]func(*Ctx)),
			userStates: NewUserStateManager(),
			log:        log.Default(),
			messages:   defaultMessages,
			ctx:        ctx,
			cancel:     cancel},

		nil
}

func (t *TgHandler) SetLogger(logger Logger) {
	t.log = logger
}

func (t *TgHandler) HandleCmd(command string, f func(*Ctx)) {
	t.croutes[command] = f
}
func (t *TgHandler) HandleState(command string, f func(*Ctx)) {
	t.sroutes[command] = f
}

func (t *TgHandler) processUpdate(ctx context.Context, update *tgbotapi.Update) {
	if update.Message != nil {
		userID := update.Message.From.ID
		userState := t.userStates.GetUserState(int64(userID))
		mctx := &Ctx{update.Message, t.bot, userState, t.log}

		currentFunction := userState.GetCurrentFunction()
		if currentFunction != "" {
			handler, ok := t.sroutes[currentFunction]
			if ok {
				go func() {
					select {
					case <-ctx.Done():
						return
					default:
						handler(mctx)
					}
				}()
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
			go func() {
				select {
				case <-ctx.Done():
					return
				default:
					handler(mctx)
				}
			}()
		} else {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, t.messages.NoCommand))
		}
	}
}

func (t *TgHandler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case <-t.ctx.Done():
			return
		case update := <-updates:
			t.processUpdate(t.ctx, &update)
		}
	}
}

func (t *TgHandler) Stop() {
	t.cancel()
}

func (t *TgHandler) SetCustomMessages(messages *Messages) {
	t.messages = messages
}
