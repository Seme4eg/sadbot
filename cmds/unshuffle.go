package cmds

import "fmt"

func Unshuffle(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream.UnShuffle()
	ctx.Reply("Shuffling turned off.")
}
