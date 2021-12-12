package plugins

import (
	"github.com/bwmarrin/discordgo"
)

type Plugin interface {
	Init(s *discordgo.Session, guildID string, appID string)
	GetHandlers() []interface{}
	GetIntents() []discordgo.Intent
}
