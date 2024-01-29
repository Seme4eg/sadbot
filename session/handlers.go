package session

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/seme4eg/sadbot/cmd"
)

// ready is a handler for the "ready" event from Discord
func (s *Session) ready() func(*discordgo.Session, *discordgo.Ready) {
	return func(ses *discordgo.Session, event *discordgo.Ready) {
		ses.UpdateListeningStatus(s.prefix + "help")
	}
}

// messageCreate is a handler for "message create" event from Discord
func (s *Session) messageCreate() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(ses *discordgo.Session, m *discordgo.MessageCreate) {
		// ignore bot messages and all messages without bot prefix
		if m.Author.ID == ses.State.User.ID || !func() bool {
			str := m.Content
			return len(str) >= len(string(s.prefix)) && str[0:len(string(s.prefix))] == string(s.prefix)
		}() {
			return
		}

		command := strings.TrimPrefix(m.Content, s.prefix)

		if len(command) < 1 {
			return
		}

		cmd.Handle(command, ses, m, s.Streams, s.prefix)
	}
}

// guildCreate is a handler to be called every time a new guild is joined
func (s *Session) guildCreate() func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(ses *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Unavailable {
			return
		}
		for _, channel := range event.Channels {
			if channel.ID == event.ID {
				_, err := ses.ChannelMessageSend(
					channel.ID,
					fmt.Sprintf("sadbot is ready! You need %shelp.", s.prefix))
				if err != nil {
					fmt.Println("Failed to greet the guild:", err)
				}
				return
			}
		}
	}
}
