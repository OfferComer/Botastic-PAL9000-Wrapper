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

func Init(ctx context.Context, name string, cfg config.MixinCo