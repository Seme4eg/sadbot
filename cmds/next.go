package cmds

import "fmt"

func Next(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ctx.Stream.Queue) == 0 {
		ctx.Reply("Nothing to play, add tracks first.")
		return
	}
	ctx.Stream.Playing = true
	if len(ctx.Stream.Queue) == 1 {
		ctx.Reply("Last track in queue")
		return
	}
	// sends stop signal to current ffmpeg command stopping it
	ctx.Stream.Stop <- true
}
