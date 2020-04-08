package app

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-vk-api/vk"
	"github.com/stretchr/testify/assert"

	"github.com/petuhovskiy/vk-watcher-bot/conf"
)

func TestVkClient(t *testing.T) {
	cfg, err := conf.ParseEnv()
	assert.Nil(t, err)

	client, err := vk.NewClientWithOptions(
		vk.WithToken(cfg.VK.Token),
	)
	assert.Nil(t, err)

	cli := NewVkClient(client)
	spew.Dump(cli.ReadTopic(vk.RequestParams{
		"group_id": "128026383",
		"topic_id": "41000758",
		"extended": 1,
		"count":    50,
	}))
}
