package app

import (
	"fmt"
	"github.com/petuhovskiy/vk-watcher-bot/conf"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/petuhovskiy/telegram"
)

type Handler struct {
	bot   *telegram.Bot
	logic *Logic
	cfg   *conf.Struct
}

func NewHandler(bot *telegram.Bot, logic *Logic, cfg *conf.Struct) *Handler {
	return &Handler{
		bot:   bot,
		logic: logic,
		cfg:   cfg,
	}
}

func (h *Handler) Handle(upd *telegram.Update) {
	if upd.Message == nil {
		return
	}

	msg := upd.Message
	if msg.From.ID != h.cfg.Bot.AdminID {
		return
	}

	h.handleMessage(msg)
}

func (h *Handler) handleMessage(msg *telegram.Message) {
	logrus.WithField("msg", msg).Info("incoming message")

	text := msg.Text
	if !strings.HasPrefix(text, "/") {
		return
	}

	cmds := strings.Split(text, " ")
	h.handleCommand(msg.Chat.ID, cmds)
}

func (h *Handler) handleCommand(chatID int, cmds []string) {
	if len(cmds) == 0 {
		return
	}

	cmd := cmds[0]
	switch cmd {
	default:
		h.commandNotFound(chatID)
	}
}

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help? Ain't got some.
`

	NewSender(h.bot, fmt.Sprintf("%d", chatID)).Send(str)
}
