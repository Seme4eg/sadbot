package cmd

import "fmt"

// Unshuffle calls current guild's Unshuffle method.
func Unshuffle(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream().UnShuffle()
	ctx.Reply("Shuffling turned off.")
}
