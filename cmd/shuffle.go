package cmd

import "fmt"

// Shuffle calls current guild's stream Shuffle method. Replies with success msg.
func Shuffle(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.stream().Shuffle()
	ctx.reply("Shuffling turned on.")
}
