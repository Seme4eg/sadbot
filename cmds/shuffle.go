package cmds

import "fmt"

// Shuffle calls current guild's stream Shuffle method. Replies with success msg.
func Shuffle(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream().Shuffle()
	ctx.Reply("Shuffling turned on.")
}
