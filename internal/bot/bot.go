package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/handler"
	"github.com/nathanim1919/mezgeb/internal/bot/state"
	"github.com/nathanim1919/mezgeb/internal/service"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *handler.Handler
}

func New(token string, svc *service.Service) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	stateMgr := state.NewManager()
	h := handler.New(api, svc, stateMgr)

	return &Bot{api: api, handler: h}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Bot is running. Waiting for messages...")

	for update := range updates {
		go b.handler.HandleUpdate(update)
	}
}
