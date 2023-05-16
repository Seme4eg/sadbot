package cmds

import "fmt"

func Next(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if errMessage := ctx.Stream.Next(); errMessage != "" {
		ctx.Reply(errMessage)
	} else {
		ctx.Reply("Now playing: " + ctx.Stream.Current())
	}
}
