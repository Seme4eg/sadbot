package cmds

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Ctx struct {
	S    *discordgo.Session
	M    *discordgo.MessageCreate
	Args []string
}

func (c Ctx) Reply(msg string) {
	_, err := c.S.ChannelMessageSend(c.M.ChannelID, msg)
	if err != nil {
		fmt.Println("Failed to send message to channel: ", err.Error())
	}
}
