package cmds

import "fmt"

func Leave(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return in case bot is already connected to some channel
	if ctx.Stream.V != nil {
		err := ctx.Stream.V.Disconnect()
		if err != nil {
			fmt.Println("Error leaving voice channel:", err)
		}
	} else {
		ctx.Reply("I'm not in vc")
		return
	}
}
