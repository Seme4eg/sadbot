package cmd

import "fmt"

// Leave disconnects bot from current voice channel (in current guild) and
// removes guild's Stream from the pool of all streams.
func Leave(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return in case bot is already connected to some channel
	if ctx.stream().V != nil {
		// TODO: here (or maybe in leave command ye), and on 'join' event play
		// some japanese sounds like 'hi' and 'bye' , better from anime
		err := ctx.stream().Disconnect()
		if err != nil {
			fmt.Println("Error leaving voice channel:", err)
			return
		}
		// delete stream of this guild from streams map
		delete(ctx.Streams.List, ctx.M.GuildID)
	} else {
		ctx.reply("You can check out any time you like, but you can never leave.")
		return
	}
}
