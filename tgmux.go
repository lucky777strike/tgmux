package tgmux

import (
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
}

func NewHandler(token string) (*TgHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TgHandler{bot: bot,
			croutes:    make(map[string]func(*Ctx)),
			sroutes:    make(map[string]func(*Ctx)),
			userStates: NewUserStateManager()},
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
		mctx := &Ctx{update.Message, t.bot, userState}

		if userState.CurrentFunction != "" {
			handler, ok := t.sroutes[userState.CurrentFunction]
			if ok {
				go handler(mctx)
			} else {
				t.userStates.ResetUserFunction((int64(userID)))
				t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error,try again"))
			}
			return
		}
		handler, ok := t.croutes[update.Message.Text]
		if ok {
			go handler(mctx)
		} else {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No such command"))
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
