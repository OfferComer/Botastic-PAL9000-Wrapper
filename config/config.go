
package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General  GeneralConfig  `yaml:"general"`
	Adapters AdaptersConfig `yaml:"adapters"`
}

func (s *Config) String() string {
	data, _ := yaml.Marshal(s)
	return string(data)
}

type BotConfig struct {
	BotID uint64 `yaml:"bot_id"`
	Lang  string `yaml:"lang"`
}

type BotasticConfig struct {
	AppId string `yaml:"app_id"`
	Host  string `yaml:"host"`
	Debug bool   `yaml:"debug"`
}

type GeneralOptionsConfig struct {
	IgnoreIfError bool `yaml:"ignore_if_error"`
	FormatLinks   bool `yaml:"format_links"`
}

type GeneralConfig struct {
	Options  *GeneralOptionsConfig `yaml:"options,omitempty"`
	Bot      *BotConfig            `yaml:"bot,omitempty"`
	Botastic *BotasticConfig       `yaml:"botastic,omitempty"`
}

type AdaptersConfig struct {
	Enabled []string                 `yaml:"enabled"`
	Items   map[string]AdapterConfig `yaml:"items"`
}

type AdapterConfig struct {
	Driver   string          `yaml:"driver"`
	Mixin    *MixinConfig    `yaml:"mixin,omitempty"`
	Telegram *TelegramConfig `yaml:"telegram,omitempty"`
	Discord  *DiscordConfig  `yaml:"discord,omitempty"`
	WeChat   *WeChatConfig   `yaml:"wechat,omitempty"`
}

type WeChatConfig struct {
	GeneralConfig `yaml:",inline"`

	Address string `yaml:"address"`
	Path    string `yaml:"path"`
	Token   string `yaml:"token"`
}

type MixinConfig struct {
	GeneralConfig `yaml:",inline"`

	Keystore               string   `yaml:"keystore"` // base64 encoded keystore (json format)
	Whitelist              []string `yaml:"whitelist"`
	MessageCacheExpiration int64    `yaml:"message_cache_expiration"`
}

type TelegramConfig struct {
	GeneralConfig `yaml:",inline"`

	Debug     bool     `yaml:"debug"`
	Token     string   `yaml:"token"`
	Whitelist []string `yaml:"whitelist"`
}

type DiscordConfig struct {
	GeneralConfig `yaml:",inline"`

	Token     string   `yaml:"token"`
	Whitelist []string `yaml:"whitelist"`
}

func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			Options: &GeneralOptionsConfig{