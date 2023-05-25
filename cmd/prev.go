package cmd

import "fmt"

// Prev calls current guild's stream Prev method. On success replies
// with current track name.
func Prev(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ctx.Stream().Prev(); err != nil {
		ctx.Reply(err.Error())
	} else {
		ctx.Reply("Now playing: " + ctx.Stream().Current())
	}
}
