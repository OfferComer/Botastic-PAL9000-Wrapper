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
	BotID