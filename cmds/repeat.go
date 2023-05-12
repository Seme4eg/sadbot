package cmds

import "fmt"

func Repeat(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	reponse := ctx.Stream.SetRepeat(ctx.Args)
	ctx.Reply(reponse)
}
