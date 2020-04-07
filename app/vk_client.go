package app

import (
	"fmt"
	"github.com/go-vk-api/vk"
	"html"
	"time"
)

type TopicResponse struct {
	Count    int
	Items    []TopicItem
	Profiles []Profile
}

type Profile struct {
	Type       string `json:"type"`
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Sex        int    `json:"sex"`
	ScreenName string `json:"screen_name"`
	Photo50    string `json:"photo_50"`
	Photo100   string `json:"photo_100"`
	Online     int    `json:"online"`
}

type TopicItem struct {
	ID          int
	FromID      int       `json:"from_id"` // user id
	Date        int       // unix time
	Time        time.Time `json:"-"`
	Text        string
	Attachments []Attach

	AwesomeText string `json:"-"`
}

type Attach map[string]interface{}

type VkClient struct {
	*vk.Client
}

func NewVkClient(cli *vk.Client) VkClient {
	return VkClient{cli}
}

func (v VkClient) ReadTopic(req vk.RequestParams) (TopicResponse, error) {
	var resp TopicResponse
	err := v.CallMethod("board.getComments", req, &resp)

	for i, item := range resp.Items {
		item.Time = time.Unix(int64(item.Date), 0)

		username := "<b>Org</b>"
		for _, prof := range resp.Profiles {
			if prof.ID == item.FromID {
				username = prof.FirstName + " " + prof.LastName
			}
		}

		link := fmt.Sprintf("https://vk.com/topic-%v_%v?post=%v", req["group_id"], req["topic_id"], item.ID)
		awesomeText := fmt.Sprintf(
			"<a href=\"%v\">%v</a>\nBy: %v\n%v",
			link,
			item.Time,
			username,
			html.EscapeString(item.Text),
		)

		if len(item.Attachments) > 0 {
			awesomeText += "\n\nHAS SOME ATTACHMENTS, CHECK IT OUT!"
		}

		item.AwesomeText = awesomeText

		resp.Items[i] = item
	}

	return resp, err
}
