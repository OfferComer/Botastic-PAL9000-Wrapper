package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pandodao/PAL9000/config"
	"github.com/pandodao/PAL9000/store"
	"github.com/pandodao/botastic-go"
	"github.com/sirupsen/logrus"
)

var (
	linkRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
)

type Adapter interface {
	GetName() string
	GetMessageChan(ctx context.Context) <-chan *Message
	HandleResult(message *Message, result *Result)
}

type Handler struct {
	cfg     config.GeneralConfig
	client  *botastic.Client
	store   store.Store
	adapter Adapter
	logger  *logrus.Entry
}

type Message struct {
	Context context.Context
	BotID   uint64
	Lang    string

	UserIdentity string
	ConvKey      string
	Content      string
	ReplyContent string

	DoneChan chan struct{}
}

type Result struct {
	ConvTurn      *botastic.ConvTurn
	Err           error
	IgnoreIfError bool
}

func NewHandler(cfg config.GeneralConfig, store store.Store, adapter Adapter) *Handler {
	client := botastic.New(cfg.Botastic.AppId, "", botastic.WithDebug(cfg.Botastic.Debug), botastic.WithHost(cfg.Botastic.Host))
	return &Handler{
		cfg:     cfg,
		client:  client,
		store:   store,
		adapter: adapter,
		logger:  logrus.WithField("adapter", fmt.Sprintf("%T", adapter)).WithField("component", "service").WithField("adapter_name", adapter.GetName()),
	}
}

func (h *Handler) Start(ctx context.Context) error {
	msgChan := h.adapter.GetMessageChan(ctx)

	for {
		select {
		case msg := <-msgChan:
			h.logger.WithField("msg", msg).Info("received message")
			if msg.BotID == 0 {
				msg.BotID = h.cfg.Bot.BotID
			}
			if msg.Lang == "" {
				msg.Lang = h.cfg.Bot.Lang
			}

			turn, err := h.handleMessage(ctx, msg)
			h.logger.WithFields(logrus.Fields{
				"turn":       turn,
				"result_err": err,
			}).Info("handled message")
			h.adapter.HandleResult(msg, &Result{
				ConvTurn:      turn,
				IgnoreIfError: h.cfg.Options.IgnoreIfError,
				Err:           err,
			})
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, m *Message) (*botastic.ConvTurn, error) {
	conv, err := h.store.GetConversationByKey(m.ConvKey)
	if err != nil {
		return nil, err
	}

	if conv == nil {
		conv, err = h.client.CreateConversation(