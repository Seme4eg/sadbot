package cmds

import (
	"fmt"
)

// XXX: should this package be handling all errors or should it just return err?

func Leave(ctx Ctx) {
	// Get the voice state for the given guild and user
	VoiceState, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)

	// if err means user is not connected to a voice channel
	if err != nil {
		ctx.Reply("Must be connected to voice channel to use bot")
		return
	}

	// return in case bot is already connected to some channel
	if c, ok := ctx.S.VoiceConnections[ctx.M.GuildID]; ok {
		// forbid user to command 'leave' if in different channel than the bot
		if c.ChannelID != VoiceState.ChannelID {
			ctx.Reply("Must be in same voice channel as bot")
			return
		}
		err := c.Disconnect()
		if err != nil {
			fmt.Println("Error leaving voice channel:", err.Error())
		}
	}
}
