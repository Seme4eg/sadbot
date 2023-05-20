package cmds

func NowPlaying(ctx Ctx) {
	if current := ctx.Stream().Current(); current == "" {
		ctx.Reply("Queue is empty, sir")
	} else {
		ctx.Reply("Now playing: " + current)
	}
}
