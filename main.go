package main

import (
	"fmt"
	"github.com/seme4eg/sadbot/cmd"
	"github.com/seme4eg/sadbot/stream"
	"github.com/seme4eg/sadbot/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	config *utils.Config
	// map of all streams that contains one stream per one guild
	Streams *stream.Streams
)

func main() {

	// parse config file
	var err error
	if config, err = utils.NewConfig(); err != nil {
		fmt.Println("Failed to parse config file", err)
		os.Exit(1)
	}

	// Create new Discord Session
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Failed to create a discord session:", err)
		return
	}

	Streams = &stream.Streams{List: make(map[string]*stream.Stream)}

	session.AddHandler(ready)         // ready events.
	session.AddHandler(messageCreate) // messageCreate events.
	session.AddHandler(guildCreate)   // guildCreate events.

	// create websocket connection with discord
	if err := session.Open(); err != nil {
		fmt.Println(err)
		return
	}
	// ensure that session will be gracefully closed whenever the function exits
	defer session.Close()

	// run until code is terminated
	fmt.Println("sadbot is now running. Press Ctrl-C to exit.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c
}

// Handler for the "ready" event from Discord
func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateListeningStatus(config.Prefix + "help")
}

// Handler for "message create" event from Discord
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages and all messages without bot prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	command := strings.TrimPrefix(m.Content, config.Prefix)

	if len(command) < 1 {
		return
	}

	cmd.Handle(command, s, m, Streams, config.Prefix)
}

// This function will be called every time a new guild is joined
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Unavailable {
		return
	}
	for _, channel := range event.Channels {
		if channel.ID == event.ID {
			_, err := s.ChannelMessageSend(
				channel.ID,
				fmt.Sprintf("sadbot is ready! You need %shelp.", config.Prefix))
			if err != nil {
				fmt.Println("Failed to greet the guild:", err)
			}
			return
		}
	}
}
