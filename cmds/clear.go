package cmds

import "fmt"

func Clear(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream().Clear()
	ctx.Reply("Queue cleared")
}
