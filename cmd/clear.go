package cmd

import "fmt"

// Clear calls current guild's stream Clear method.
func Clear(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream().Clear()
	ctx.Reply("Queue cleared")
}
