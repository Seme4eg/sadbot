package cmd

// NowPlaying replies with currently playing song name if there is one.
func NowPlaying(ctx Ctx) {
	if current := ctx.Stream().Current(); current == "" {
		ctx.Reply("Queue is empty, sir")
	} else {
		ctx.Reply("Now playing: " + current)
	}
}
