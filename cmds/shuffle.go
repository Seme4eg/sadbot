package cmds

import "fmt"

func Shuffle(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream().Shuffle()
	ctx.Reply("Shuffling turned on.")
}
