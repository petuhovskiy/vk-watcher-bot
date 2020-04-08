package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-vk-api/vk"
	"github.com/petuhovskiy/vk-watcher-bot/app"
	"github.com/petuhovskiy/vk-watcher-bot/conf"

	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
)

const (
	watcherStartDelay = time.Second * 30
)

func main() {
	cfg, err := conf.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	client, err := vk.NewClientWithOptions(
		vk.WithToken(cfg.VK.Token),
	)
	if err != nil {
		log.Fatal(err)
	}

	cli := app.NewVkClient(client)

	bot := telegram.NewBotWithOpts(cfg.Bot.Token, &telegram.Opts{
		Middleware: func(handler telegram.RequestHandler) telegram.RequestHandler {
			return func(methodName string, req interface{}) (message json.RawMessage, err error) {
				res, err := handler(methodName, req)
				if err != nil {
					log.Println("Telegram response error: ", err)
				}

				return res, err
			}
		},
	})

	ch, err := updates.StartPolling(bot, telegram.GetUpdatesRequest{
		Offset:  0,
		Limit:   50,
		Timeout: 60,
	})
	if err != nil {
		log.Fatal(err)
	}

	sender := app.NewSender(bot, cfg.Bot.ChannelID)

	// TODO: fix hardcore later
	tmp := 5903
	discussWatcher := app.NewWatcher(cli, "128026383", "41000758", sender, cfg.PrevDuration, &tmp)
	go discussWatcher.Start()

	hintWatcher := app.NewWatcher(cli, "128026383", "40182106", sender, cfg.PrevDurationHint, nil)
	go func() {
		time.Sleep(watcherStartDelay)
		hintWatcher.Start()
	}()

	logic := app.NewLogic()
	handler := app.NewHandler(bot, logic, cfg)

	for upd := range ch {
		upd := upd
		handler.Handle(&upd)
	}
}
