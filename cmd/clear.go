package cmd

import "fmt"

// Clear calls current guild's stream Clear method.
func Clear(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.stream().Clear()
	ctx.reply("Queue cleared")
}
