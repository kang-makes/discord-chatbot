package bot

import (
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"github.com/kang-makes/discord-chatbot/bot/plugins"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

type Bot struct {
	Auth      BotAuth       `json:"auth" yaml:"auth"`
	Instances []BotInstance `json:"instances" yaml:"instances"`
}

type BotAuth struct {
	AppID   string `json:"appID" yaml:"appID"`
	Token   string `json:"token" yaml:"token"`
	Session *discordgo.Session
}

type BotInstance struct {
	GuildID string        `json:"guildID" yaml:"guildID"`
	Plugins []PluginEntry `json:"plugins" yaml:"plugins"`
}

type PluginEntry struct {
	Name   string    `json:"name" yaml:"name"`
	Params yaml.Node `yaml:"params"`
	plugin plugins.Plugin
}

func ParseConfig(s string) (b Bot, err error) {
	log.Trace("reading config file")
	data, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Debugf("config file read: %q", data)

	log.Trace("unmarchalling yaml")
	err = yaml.Unmarshal(data, &b)
	if err != nil {
		log.Fatal(err)
		return Bot{}, err
	}
	log.Debugf("config loaded: %v", &b)

	return
}
