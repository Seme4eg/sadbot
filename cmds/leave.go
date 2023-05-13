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
		// TODO: here (or maybe in leave command ye), and on 'join' event play
		// some japanese sounds like 'hi' and 'bye' , better from anime
		ctx.Stream.Reset(false)
		err := ctx.Stream.V.Disconnect()
		if err != nil {
			fmt.Println("Error leaving voice channel:", err)
		}
	} else {
		ctx.Reply("You can check out any time you like, but you can never leave.")
		return
	}
}
