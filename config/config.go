package config

import (
	"encoding/json"
	"errors"
	"github.com/praveenprem/testbed-slack-bot/api"
	"github.com/praveenprem/testbed-slack-bot/slack"
	"io/ioutil"
	"log"
	"os"
)

type (
	Config struct {
		api.Server
		Slack slack.Config `json:"slack"`
	}
	Loader interface {
		Load(path string) error
	}
)

func (c *Config) Load(path string) error {
	log.Printf("loading appliaction configguration")
	if !isFile(path) {
		return errors.New("configuration file doesn't exist")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	rawCfg, rErr := ioutil.ReadAll(file)
	if rErr != nil {
		return rErr
	}

	log.Printf("config loaded")

	return json.Unmarshal(rawCfg, c)
}

func isFile(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		log.Printf(err.Error())
	}
	return false
}
