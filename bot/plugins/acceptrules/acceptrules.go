package acceptrules

import (
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

type AcceptRules struct {
	ChannelID    string `json:"channelID" yaml:"channelID"`
	MessageID    string `json:"messageID" yaml:"messageID"`
	EmojiName    string `json:"emojiName" yaml:"emojiName"`
	RoleToAssign string `json:"roleToAssign" yaml:"roleToAssign"`
	GuildID      string
}

func (ar *AcceptRules) Init(s *discordgo.Session, guildID string, appID string) {
	ar.GuildID = guildID
}

func (ar *AcceptRules) GetHandlers() []interface{} {
	return []interface{}{
		ar.OnReactionAdd,
		ar.OnReactionRemove,
	}
}

func (ar *AcceptRules) OnReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.GuildID != ar.GuildID ||
		r.ChannelID != ar.ChannelID ||
		r.MessageID != ar.MessageID ||
		r.Emoji.Name != ar.EmojiName {
		log.Trace("AcceptRules.OnReactionAdd ignored. Mismatch.")
	} else {
		log.Trace("AcceptRules.OnReactionAdd matched!")
		member, _ := s.GuildMember(ar.GuildID, r.UserID)
		log.Tracef("roles from the member: %v", member.Roles)
		for _, role := range member.Roles {
			if role == ar.RoleToAssign {
				return
			}
		}
		roles := append(member.Roles, ar.RoleToAssign)
		log.Debugf("assigning roles at '%s' to '%s': %v", ar.GuildID, member.User.ID, roles)
		err := s.GuildMemberEdit(ar.GuildID, r.UserID, roles)
		if err != nil {
			user := member.User.String()
			guild, _ := s.Guild(ar.GuildID)
			log.Errorf("error while removing roles to '%s' at '%s' while trying to set %v: %v", user, guild.Name, roles, err)
		}
	}
}

func (ar *AcceptRules) OnReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.GuildID != ar.GuildID ||
		r.ChannelID != ar.ChannelID ||
		r.MessageID != ar.MessageID ||
		r.Emoji.Name != ar.EmojiName {
		log.Trace("AcceptRules.OnReactionAdd ignored. Mismatch.")
	} else {
		log.Trace("AcceptRules.OnReactionRemove matched!")
		member, _ := s.GuildMember(ar.GuildID, r.UserID)
		log.Tracef("roles from the member: %v", member.Roles)
		for i, role := range member.Roles {
			if role == ar.RoleToAssign {
				roles := append(member.Roles[:i], member.Roles[i+1:]...)
				log.Debugf("assigning roles at '%s' to '%s': %v", ar.GuildID, member.User.ID, roles)
				s.GuildMemberEdit(ar.GuildID, r.UserID, roles)
			}
		}
	}
}

func (ar *AcceptRules) GetIntents() []discordgo.Intent {
	return []discordgo.Intent{
		discordgo.IntentsGuildMessageReactions,
	}
}
