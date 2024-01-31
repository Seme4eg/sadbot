package cmd

// NowPlaying replies with currently playing track name if there is one.
func NowPlaying(ctx Ctx) {
	if current := ctx.stream().Current(); current == "" {
		ctx.reply("Queue is empty, sir")
	} else {
		ctx.reply("Now playing: " + current)
	}
}
