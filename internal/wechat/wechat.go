
package wechat

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/pandodao/PAL9000/config"
	"github.com/pandodao/PAL9000/service"
)

type httpRequsetKey struct{}
type httpResponseKey struct{}
type rawMessageKey struct{}

type TextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

type Bot struct {
	name string
	cfg  config.WeChatConfig
}

func New(name string, cfg config.WeChatConfig) *Bot {
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
	go func() {
		server := &http.Server{
			Addr:    b.cfg.Address,
			Handler: http.DefaultServeMux,
		}

		validateSignature := func(signature, timestamp, nonce string) bool {
			params := []string{b.cfg.Token, timestamp, nonce}
			sort.Strings(params)
			combined := strings.Join(params, "")

			hash := sha1.New()
			hash.Write([]byte(combined))
			hashStr := hex.EncodeToString(hash.Sum(nil))

			return hashStr == signature
		}

		http.HandleFunc(b.cfg.Path, func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			signature := r.Form.Get("signature")
			timestamp := r.Form.Get("timestamp")
			nonce := r.Form.Get("nonce")
			echostr := r.Form.Get("echostr")

			if !validateSignature(signature, timestamp, nonce) {
				http.Error(w, "Invalid signature", http.StatusForbidden)
				return
			}