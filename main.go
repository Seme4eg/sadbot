package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sadbot/cmds"
	"sadbot/stream"
	"sadbot/utils"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	config  utils.Config
	session *discordgo.Session
	Stream  *stream.Stream
)

func init() {
	if err := utils.ReadConfig(&config); err != nil {
		log.Fatal("Failed parse config file", err)
	}
}

func main() {
	var err error
	// Create new Discord Session
	if session, err = discordgo.New("Bot " + config.Token); err != nil {
		fmt.Println(err)
		return
	}

	Stream = &stream.Stream{
		S:      session,
		Stop:   make(chan bool),
		Repeat: stream.RepeatOff,
	}

	// Register ready as a callback for the ready events.
	session.AddHandler(ready)
	// Register messageCreate as a callback for the messageCreate events.
	session.AddHandler(messageCreate)
	// Register guildCreate as a callback for the guildCreate events.
	// session.AddHandler(guildCreate)

	if err := session.Open(); err != nil {
		fmt.Println(err)
		return
	}
	// ensure that session will be gracefully closed whenever the function exits
	defer session.Close()

	fmt.Println("Bot is running !")

	// run until code is terminated
	fmt.Println("sadbot is now running. Press CTRL-C to exit.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c
}

// This function will be called when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateListeningStatus(config.Prefix + "help")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages and all messages without needed prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	command := strings.TrimPrefix(m.Content, config.Prefix)

	if len(command) < 1 {
		return
	}

	name := strings.Fields(command)[0]
	args := strings.TrimPrefix(command, name+" ")
	if command, ok := cmds.Pool[strings.ToLower(name)]; ok {
		ctx := cmds.Ctx{S: s, M: m, Args: args, Stream: Stream}
		command(ctx)
	}
}

// This function will be called every time a new guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Unavailable {
		return
	}
	for _, channel := range event.Channels {
		if channel.ID == event.ID {
			_, _ = s.ChannelMessageSend(
				channel.ID,
				fmt.Sprintf(
					"sadbot is ready! Type %shelp to see what it's capable of.",
					config.Prefix))
			return
		}
	}
}
