package cmd

import "fmt"

// Unshuffle calls current guild's Unshuffle method.
func Unshuffle(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.stream().UnShuffle()
	ctx.reply("Shuffling turned off.")
}
