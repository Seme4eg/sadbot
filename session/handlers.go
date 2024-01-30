package session

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/seme4eg/sadbot/cmd"
)

// ready is a handler for the "ready" event from Discord
func (s *Session) ready(ses *discordgo.Session, event *discordgo.Ready) {
	ses.UpdateListeningStatus(s.prefix + "help")
}

// messageCreate is a handler for "message create" event from Discord
func (s *Session) messageCreate(ses *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages and all messages without bot prefix
	if m.Author.ID == ses.State.User.ID || !func() bool {
		str := m.Content
		return len(str) >= len(string(s.prefix)) && str[0:len(string(s.prefix))] == string(s.prefix)
	}() {
		return
	}

	cmd.Handle(ses, m, s.Streams, s.prefix)
}

// guildCreate is a handler to be called every time a new guild is joined
func (s *Session) guildCreate(ses *discordgo.Session, event *discordgo.GuildCreate) {
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
