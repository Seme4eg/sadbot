package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Start(session *discordgo.Session, config Config) (err error) {
	// Create new Discord Session
	if session, err = discordgo.New("Bot " + config.Token); err != nil {
		return
	}

	session.AddHandler(messageHandler) // adding event handler

	if err = session.Open(); err != nil {
		return
	}
	// ensure that session will be gracefully closed whenever the function exits
	defer session.Close()

	fmt.Println("Bot is running !")

	// run until code is terminated
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages and all messages without needed prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	command := strings.TrimPrefix(m.Content, config.Prefix)

	switch command {
	case "ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
