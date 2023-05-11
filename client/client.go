package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/expose443/BookLinkerBot/flibusta"
	"github.com/expose443/BookLinkerBot/model"
)

type MessageService struct {
	*http.Client
	botUrl string
	Offset int
}

func NewHttpClient(botUrl string) *MessageService {
	return &MessageService{
		&http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Println(req.Response.Status)
				fmt.Println("[REDIRECT]")
				return nil
			},
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 30,
		},
		botUrl,
		0,
	}
}

func (m *MessageService) GetUpdates() ([]model.Update, error) {
	resp, err := http.Get(m.botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(m.Offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse model.RestResponse
	if err = json.Unmarshal(body, &restResponse); err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func (m *MessageService) Respond(update model.Update) error {
	fmt.Println(update.Info())

	var botMessage model.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId

	switch update.Message.Text {
	case "/start":
		botMessage.Text = "Write book name: "
	default:
		botMessage.Text = flibusta.GetBookLinks(update.Message.Text)
	}

	botMessage.Text = flibusta.GetBookLinks(strings.ReplaceAll(update.Message.Text, " ", "+"))

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	if _, err = m.Post(m.botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf)); err != nil {
		return err
	}
	return nil
}
