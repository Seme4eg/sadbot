package cmds

import "fmt"

func Shuffle(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch ctx.Args {
	case "off":
		ctx.Stream.Shuffle()
		ctx.Reply("Shuffling turned off.")
	default:
		ctx.Stream.UnShuffle()
		ctx.Reply("Shuffling turned on.")
	}
}
