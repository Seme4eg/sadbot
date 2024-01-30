package cmd

import "fmt"

// Prev calls current guild's stream Prev method. On success replies
// with current track name.
func Prev(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ctx.stream().Prev(); err != nil {
		ctx.reply(err.Error())
	} else {
		ctx.reply("Now playing: " + ctx.stream().Current())
	}
}
