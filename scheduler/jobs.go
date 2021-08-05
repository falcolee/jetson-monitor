package scheduler

import (
	"github.com/juandiii/jetson-monitor/config"
	"github.com/juandiii/jetson-monitor/logging"
	"github.com/juandiii/jetson-monitor/notification"
	"github.com/juandiii/jetson-monitor/notification/dingding"
	"github.com/juandiii/jetson-monitor/notification/slack"
	"github.com/juandiii/jetson-monitor/notification/telegram"
	"github.com/juandiii/jetson-monitor/request"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"time"
)

type Scheduler struct {
	conf                  config.URL
	NotificationProviders []notification.CommandProvider
	Logger                *logging.StandardLogger
	Cache                 *cache.Cache
}

func New(c config.URL, conf *config.ConfigJetson) cron.Job {
	return &Scheduler{
		conf: c,
		NotificationProviders: []notification.CommandProvider{
			slack.New(c, conf.Logger),
			telegram.New(c, conf.Logger),
			dingding.New(c, conf.Logger),
		},
		Logger: conf.Logger,
		Cache:  conf.Cache,
	}
}

func (s *Scheduler) Run() {
	_, err := request.RequestServer(s.conf, s.Logger)
	if err != nil {
		_, found := s.Cache.Get(s.conf.URL)
		if !found {
			var notifyErr error
			for _, n := range s.NotificationProviders {
				notifyErr = n.SendMessage(&notification.Message{
					Text: err.Error(),
				})
			}
			if notifyErr == nil {
				if s.conf.NotifyInterval != nil {
					s.Cache.Set(s.conf.URL, 1, time.Duration(*s.conf.NotifyInterval)*time.Minute)
				} else {
					s.Cache.Set(s.conf.URL, 1, cache.DefaultExpiration)
				}
			}
		}
	}
}
