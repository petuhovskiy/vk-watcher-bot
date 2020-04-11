package app

import (
	"time"

	"github.com/go-vk-api/vk"
	"github.com/sirupsen/logrus"
)

const (
	readTopicDelay = time.Second * 3
)

type Watcher struct {
	cli     VkClient
	groupID string
	topicID string
	sender  *Sender
	dur     time.Duration

	skipAnyway     time.Time
	alreadySent    map[string]struct{}
	startCommentID *int

	messagePin bool // If it equals to true, we should pin all the sent messages
}

func NewWatcher(cli VkClient, groupID, topicID string, sender *Sender, dur time.Duration, start *int, messagePin bool) *Watcher {
	return &Watcher{
		cli:            cli,
		groupID:        groupID,
		topicID:        topicID,
		sender:         sender,
		dur:            dur,
		skipAnyway:     time.Now().Add(-dur),
		alreadySent:    map[string]struct{}{},
		startCommentID: start,
		messagePin:     messagePin,
	}
}

func (w *Watcher) Start() {
	for {
		w.readAll()

		// TODO: maybe ticker?
		time.Sleep(time.Minute)
	}
}

func (w *Watcher) readAll() {
	req := vk.RequestParams{
		"group_id": w.groupID,
		"topic_id": w.topicID,
		"extended": 1,
		"count":    50,
	}

	if w.startCommentID != nil {
		req["start_comment_id"] = *w.startCommentID
	}

	for offset := 0; ; offset += 50 {
		req["offset"] = offset

		resp, err := w.cli.ReadTopic(req)
		logrus.WithField("resp", resp).WithField("req", req).Info("req resp")
		if err != nil {
			logrus.WithError(err).Error("read topic error")
			break
		}
		time.Sleep(readTopicDelay)

		if len(resp.Items) == 0 {
			break
		}

		for _, item := range resp.Items {
			if time.Since(item.Time) >= w.dur {
				logrus.WithField("comment_id", item.ID).WithField("since", time.Since(item.Time)).Info("updated comment id")
				w.startCommentID = &item.ID
			}

			_, ok := w.alreadySent[item.AwesomeText]
			if ok {
				continue
			}

			w.alreadySent[item.AwesomeText] = struct{}{}

			if item.Time.Before(w.skipAnyway) {
				logrus.WithField("text", item.AwesomeText).Info("skipping anyway old post")
				continue
			}

			logrus.WithField("text", item.AwesomeText).Info("send content")
			msg, err := w.sender.Send(item.AwesomeText)
			if w.messagePin && (err == nil) {
				_ = w.sender.Pin(msg)
			}
		}
	}
}
