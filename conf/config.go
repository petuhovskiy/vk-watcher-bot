package conf

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type Struct struct {
	Bot              Bot
	VK               VK
	PrevDuration     time.Duration `env:"PREV_DURATION" envDefault:"200h"`
	PrevDurationHint time.Duration `env:"PREV_DURATION_HINT" envDefault:"200h"`
}

type Bot struct {
	AdminID   int    `env:"ADMIN_TELEGRAM_ID"`
	ChannelID string    `env:"CHANNEL_ID"`
	Token     string `env:"BOT_TOKEN"`
}

type VK struct {
	Token string `env:"VK_TOKEN"`
}

func ParseEnv() (*Struct, error) {
	cfg := Struct{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
