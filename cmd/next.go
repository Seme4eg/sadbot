package cmd

import "fmt"

// Next calls current guild's stream Next method. On success replies
// with current track name.
func Next(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ctx.Stream().Next(); err != nil {
		ctx.Reply(err.Error())
	} else {
		ctx.Reply("Now playing: " + ctx.Stream().Current())
	}
}
