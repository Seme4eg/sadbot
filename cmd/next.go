package cmd

import "fmt"

// Next calls current guild's stream Next method. On success replies
// with current track name.
func Next(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ctx.stream().Next(); err != nil {
		ctx.reply(err.Error())
	} else {
		ctx.reply("Now playing: " + ctx.stream().Current())
	}
}
