package app

import (
	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"

	"time"
)

const (
	telegramQueriesDelay = time.Second * 5
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

// Send sends message to the channel
func (s *Sender) Send(text string) (*telegram.Message, error) {
	msg, err := s.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID:                s.chatID,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
		Text:                  text,
	})
	if err != nil {
		log.WithError(err).Error("unable to send the message")
	}

	// TODO: workaround telegram spam ban
	time.Sleep(telegramQueriesDelay)
	return msg, err
}

func (s *Sender) Pin(msg *telegram.Message) error {
	_, err := s.bot.PinChatMessage(&telegram.PinChatMessageRequest{
		ChatID:              s.chatID,
		MessageID:           msg.MessageID,
		DisableNotification: true,
	})
	if err != nil {
		log.WithField("message", msg).WithError(err).Error("unable to pin the message")
	}

	// TODO: workaround telegram spam ban
	time.Sleep(telegramQueriesDelay)
	return err
}

func (s *Sender) SendPhoto(name string, b []byte) error {
	_, err := s.bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID: s.chatID,
		Photo:  NewBytesUploader(name, b),
	})
	return err
}
