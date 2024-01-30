package cmd

import (
	"fmt"
)

// Stop calls current guild's stream Reset method.
func Stop(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx.stream().Reset()
	ctx.reply("Player stopped, queue cleared")
}
