package cmds

var Pool = map[string]func(ctx Ctx){
	"ping": func(ctx Ctx) {
		_, _ = ctx.S.ChannelMessageSend(ctx.M.ChannelID, "pong")
	},
	"leave": Leave,
	"play":  Play,
	"p":     Play,
}
