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

func (b *Bot) GetMessageChan(ctx context.Context) <-chan *service.Message {
	msgChan := make(chan *service.Message)
	go func() {
		u := tgbotapi.NewUpdate(0)
		updates := b.client.GetUpdatesChan(u)
		for update := range updates {
			if update.Message == nil || update.Message.Chat == nil || update.Message.Text == "" {
				continue
			}

			allowed := len(b.cfg.Whitelist) == 0
			for _, id := range b.cfg.Whitelist {
				if strconv.FormatInt(update.Message.Chat.