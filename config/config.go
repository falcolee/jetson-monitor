package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juandiii/jetson-monitor/logging"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

type URL struct {
	URL            string `yaml:"url"`
	Match          string `yaml:"match"`
	Timeout        int64  `yaml:"timeout"`
	ResponseTime   *int   `yaml:"response_time"`
	NotifyInterval *int64 `yaml:"notify_interval"`
	StatusCode     *int   `yaml:"status_code"`
	SlackToken     string `yaml:"slack_token"`
	TelegramToken  string `yaml:"telegram_token"`
	DingdingToken  string `yaml:"dingding_token"`
	DingdingTitle  string `yaml:"dingding_title"`
	Scheduler      string `yaml:"scheduler"`
}

type ConfigJetson struct {
	Urls []URL `yaml:"urls"`

	Logger *logging.StandardLogger
	Cache  *cache.Cache
}

//Load Configuration
func (c *ConfigJetson) LoadConfig() (*ConfigJetson, error) {

	// log := logging.Logger

	err := ValidatePath("config.yml")

	if err != nil {
		c.Logger.Error("failed load config.yml")
		return nil, err
	}

	file, err := os.Open(filepath.Clean("config.yml"))

	if err != nil {
		return nil, err
	}

	defer file.Close()

	ymlFile := yaml.NewDecoder(file)

	if err := ymlFile.Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}

func ValidatePath(path string) error {
	// Check path if exists
	s, err := os.Stat(path)

	if err != nil {
		return err
	}

	// Check is directory

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory", path)
	}

	return nil
}
