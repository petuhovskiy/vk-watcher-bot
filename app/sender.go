package app

import (
	"github.com/petuhovskiy/telegram"
	"time"
)

type Sender struct {
	bot    *telegram.Bot
	chatID string
}

func NewSender(bot *telegram.Bot, chatID string) *Sender {
	return &Sender{
		bot:    bot,
		chatID: chatID,
	}
}

func (s *Sender) Send(text string) {
	_, _ = s.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID:                s.chatID,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
		Text:                  text,
	})

	// TODO: workaround telegram spam ban
	time.Sleep(time.Second * 5)
}

func (s *Sender) SendPhoto(name string, b []byte) error {
	_, err := s.bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID: s.chatID,
		Photo:  NewBytesUploader(name, b),
	})
	return err
}
