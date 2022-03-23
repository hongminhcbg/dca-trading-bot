package noti

import (
	"dca-bot/conf"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type TelegramNoti interface {
	Send(string) error
}

func NewTelegramNoti(cfg *conf.Config) (TelegramNoti, error) {
	if len(cfg.Noti.Url) == 0 || len(cfg.Noti.ChatId) == 0 {
		return nil, fmt.Errorf("missing config")
	}

	return &notiImpl{
		chatId: cfg.Noti.ChatId,
		url:    cfg.Noti.Url,
	}, nil
}

type notiImpl struct {
	chatId string
	url    string
}

func (n *notiImpl) Send(text string) error {
	data := url.Values{}
	data.Set("chat_id", n.chatId)
	data.Set("text", text)

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, n.url, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status response id %d, not OK", res.StatusCode)
	}

	return nil
}
