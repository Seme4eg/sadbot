package cmds

import (
	"fmt"
)

func Stop(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx.Stream().Reset()
	ctx.Reply("Player stopped, queue cleared")
}
