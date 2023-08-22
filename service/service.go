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
	GetName() 