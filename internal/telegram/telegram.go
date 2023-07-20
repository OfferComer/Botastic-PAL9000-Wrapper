package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pandodao/PAL9000/config"
	"github.com/pandodao/PAL9000/service"
)

var _ service.Adapter = (*Bot)(nil)

type (
	messageKey struct{}
)

type Bot struct {
	name   string
	cfg    config.TelegramConfig
	client *tgbotapi.BotAPI
}

func Init(name string, cfg config.TelegramConfig) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = cfg.Debug

	return &Bot{
		name:   name,
		cfg:    cfg,
		client: bot,
	}, nil
}

func (b *Bot) GetName() string {
	return b.name
}

func (b *Bot) GetMessageChan(ctx context.Context) <-chan *service.Messag