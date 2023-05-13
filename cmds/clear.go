package cmds

import "fmt"

func Clear(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream.Queue = ctx.Stream.Queue[:0]
	ctx.Stream.SongIndex = 0
	ctx.Reply("Queue cleared")
}
