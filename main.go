package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sadbot/cmds"
	"sadbot/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	config  utils.Config
	session *discordgo.Session
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
		fmt.Println(err.Error())
		return
	}

	session.AddHandler(CommandHandler) // adding event handler

	if err := session.Open(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// ensure that session will be gracefully closed whenever the function exits
	defer session.Close()

	fmt.Println("Bot is running !")

	// run until code is terminated
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages and all messages without needed prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	command := strings.TrimPrefix(m.Content, config.Prefix)

	if len(command) < 1 {
		return
	}

	args := strings.Fields(command)
	name := strings.ToLower(args[0])
	if command, ok := cmds.Pool[name]; ok {
		ctx := cmds.Ctx{
			S:    s,
			M:    m,
			Args: args[1:], // strip command itself
		}
		command(ctx)
	}
}
