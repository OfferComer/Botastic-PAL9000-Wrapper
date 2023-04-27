
package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pandodao/PAL9000/config"
	"github.com/pandodao/PAL9000/service"
)

var _ service.Adapter = (*Bot)(nil)

type messageKey struct{}
type sessionKey struct{}

type Bot struct {
	name string
	cfg  config.DiscordConfig
}

func New(name string, cfg config.DiscordConfig) *Bot {
	return &Bot{
		name: name,
		cfg:  cfg,
	}
}

func (b *Bot) GetName() string {
	return b.name
}

func (b *Bot) GetMessageChan(ctx context.Context) <-chan *service.Message {
	msgChan := make(chan *service.Message)

	dg, _ := discordgo.New("Bot " + b.cfg.Token)
	dg.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		// only text message
		if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
			return
		}

		allowed := len(b.cfg.Whitelist) == 0
		for _, id := range b.cfg.Whitelist {
			if id == m.Author.ID || (m.GuildID != "" && id == m.GuildID) {
				allowed = true
				break
			}
		}
