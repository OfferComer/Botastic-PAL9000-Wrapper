package mixin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	"github.com/pandodao/PAL9000/config"
	"github.com/pandodao/PAL9000/service"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Content string
	UserID  string
}

type (
	messageKey struct{}
	userKey    struct{}
	convKey    struct{}
)

var _ service.Adapter = (*Bot)(nil)

type Bot struct {
	name    string
	convMap map[string]*mixin.Conversation
	userMap map[string]*mixin.User

	client       *mixin.Client
	msgChan      chan *service.Message
	me           *mixin.User
	cfg          config.MixinConfig
	logger       logrus.FieldLogger
	messageCache *cache.Cache
}

func Init(ctx context.Context, name string, cfg config.MixinConfig) (*Bot, error) {
	data, err := base64.StdEncoding.DecodeString(cfg.Keystore)
	if err != nil {
		return nil, fmt.Errorf("base64 decode keystore error: %w", err)
	}

	var keystore mixin.Keystore
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, fmt.Errorf("json unmarshal keystore error: %w", err)
	}

	client, err := mixin.NewFromKeystore(&keystore)
	if err != nil {
		return nil, fmt.Errorf("mixin.NewFromKeystore error: %w", err)
	}

	me, err := client.UserMe(ctx)
	if err != nil {
		return nil, fmt.Errorf("mixinClient.UserMe error: %w", err)
	}

	if cfg.MessageCacheExpiration == 0 {
		cfg.MessageCacheExpiration = 60 * 60 * 24
	}

	return &Bot{
		name:         name,
		convMap:      make(map[string]*mixin.Conversation),
		userMap:      make(map[string]*mixin.User),
		client:       client,
		msgChan:      make(chan *service.Message),
		cfg:          cfg,
		me:           me,
		logger:       logrus.WithField("adapter", "mixin").WithField("name", name),
		messageCache: cache.New(time.Duration(cfg.MessageCacheExpiration)*time.Second, 10*time.Minute),
	}, nil
}

func (b *Bot) GetName() string {
	return b.name
}

func (b *Bot) GetMessageChan(ctx context.Context) <-chan *service.Message {
	go func() {
		for {