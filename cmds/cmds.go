package cmds

import "github.com/bwmarrin/discordgo"

type Ctx struct {
	S *discordgo.Session
	M *discordgo.MessageCreate
}

var Pool = map[string]func(ctx Ctx){
	"ping": func(ctx Ctx) {
		_, _ = ctx.S.ChannelMessageSend(ctx.M.ChannelID, "pong")
	},
}
