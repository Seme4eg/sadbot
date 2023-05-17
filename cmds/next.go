package cmds

import "fmt"

func Next(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ctx.Stream.Next(); err != nil {
		ctx.Reply(err.Error())
	} else {
		ctx.Reply("Now playing: " + ctx.Stream.Current())
	}
}
