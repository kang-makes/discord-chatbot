package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/kang-makes/discord-chatbot/bot/plugins/acceptrules"

	log "github.com/sirupsen/logrus"
)

func (b *Bot) Run(sig chan os.Signal) (err error) {
	log.Trace("creating session")
	session, err := discordgo.New(fmt.Sprintf("Bot %s", b.Auth.Token))
	if err != nil {
		return
	}

	log.Debug("session created")
	b.Auth.Session = session

	log.Trace("creating handlers and intents lists")
	var intents discordgo.Intent

	for _, guild := range b.Instances {
		log.Trace("configuring plugins for guild ", guild.GuildID)
		for _, p := range guild.Plugins {
			switch p.Name {
			case "acceptrules":
				log.Trace("initializing 'acceptrules' plugin")
				ar := &acceptrules.AcceptRules{}
				p.Params.Decode(ar)
				p.plugin = ar
			default:
				return fmt.Errorf("plugin does not exist: %s", p.Name)
			}

			p.plugin.Init(b.Auth.Session, guild.GuildID, b.Auth.AppID)

			pHandlers := p.plugin.GetHandlers()
			for i, v := range pHandlers {
				log.Tracef("installing handler from '%s' with index %d", p.Name, i)
				b.Auth.Session.AddHandler(v)
			}

			pIntents := p.plugin.GetIntents()
			for _, intent := range pIntents {
				intents |= intent
			}

			log.Debugf("'%s' for GID %s initialized: %v", p.Name, guild.GuildID, p.plugin)
		}
	}

	b.Auth.Session.Identify.Intents = intents

	err = b.Auth.Session.Open()
	if err != nil {
		return fmt.Errorf("discord session could not be opened: %s", err.Error())
	}

	log.Info("bot listening to events")
	<-sig

	log.Info("closing gratefully")
	err = b.Auth.Session.Close()
	return
}
