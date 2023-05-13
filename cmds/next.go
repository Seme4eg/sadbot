package cmds

import "fmt"

func Next(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	errMessage := ctx.Stream.Next()
	if errMessage != "" {
		ctx.Reply(errMessage)
	} else {
		// TODO: maybe if there gonna b many times i will need 'current' song
		// make it a 'Stream' method
		// for now it's here and in 'nowplaying' command
		ctx.Reply("Now playing: " + ctx.Stream.Queue[ctx.Stream.SongIndex].Title)
	}
}
