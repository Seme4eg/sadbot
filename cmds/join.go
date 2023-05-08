package cmds

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Join(ctx Ctx) (*discordgo.VoiceConnection, error) {
	// Get the voice state for the given guild and user
	VoiceState, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)

	// if err means user is not connected to a voice channel
	if err != nil {
		ctx.Reply("Must be connected to voice channel to use bot")
		return nil, err
	}

	// return in case bot is already connected to some channel
	if c, ok := ctx.S.VoiceConnections[ctx.M.GuildID]; ok {
		return c, nil
	}

	c, err := ctx.S.ChannelVoiceJoin(ctx.M.GuildID, VoiceState.ChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice: ", err.Error())
		return nil, err
	}
	return c, nil
}
