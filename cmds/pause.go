package cmds

import "fmt"

func Pause(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream.Playing = false
}
