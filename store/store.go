package store

import (
	"sync"

	"github.com/pandodao/botastic-go"
)

type Store interface {
	GetConversationByKey(key string) (*botastic.Conversat