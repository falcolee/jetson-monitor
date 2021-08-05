package dingding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/juandiii/jetson-monitor/config"
	"github.com/juandiii/jetson-monitor/logging"
	"github.com/juandiii/jetson-monitor/notification"
)

type Dingding struct {
	httpClient    http.Client
	URL           string
	DingdingToken string
	DingdingTitle string
	Logger        *logging.StandardLogger
}

func New(c config.URL, log *logging.StandardLogger) notification.CommandProvider {
	return &Dingding{
		httpClient: http.Client{
			Timeout: time.Duration(15 * time.Second),
		},
		URL:           "https://oapi.dingtalk.com/robot/send?access_token=" + c.DingdingToken,
		DingdingToken: c.DingdingToken,
		DingdingTitle: c.DingdingTitle,
		Logger:        log,
	}
}

func (s *Dingding) SendMessage(data *notification.Message) error {

	log := s.Logger

	if s.DingdingToken != "" && s.DingdingTitle != "" {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(&Message{
			MsgType: "markdown",
			Markdown: MarkdownBody{
				Title: s.DingdingTitle,
				Text:  data.Text,
			},
		})

		req, _ := http.NewRequest("POST", s.URL, buf)
		req.Header.Set("Content-Type", "application/json")

		if data != nil {
			log.Debug("Sending Message to Dingding")
			res, e := s.httpClient.Do(req)

			if e != nil {
				return e
			}

			defer res.Body.Close()
			bodyString := ""
			if res.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return err
				}
				bodyString = string(bodyBytes)
				log.Debugf("Dingding notify response:%v", bodyString)
			} else {
				return fmt.Errorf("Dingding notify failed,http error code:%d", res.StatusCode)
			}
		}
	}

	return nil
}
