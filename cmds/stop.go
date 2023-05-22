package cmds

import (
	"fmt"
)

// Stop calls current guild's stream Reset method.
func Stop(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx.Stream().Reset()
	ctx.Reply("Player stopped, queue cleared")
}
