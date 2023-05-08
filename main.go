package main

import (
	"fmt"
	"sadbot/bot"

	"github.com/bwmarrin/discordgo"
)

var (
	config  *bot.Config
	session *discordgo.Session
)

func main() {
	if err := bot.ReadConfig(config); err != nil {
		return
	}

	err := bot.Start(session, *config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
