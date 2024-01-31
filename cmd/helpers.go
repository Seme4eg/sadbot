package cmd

// Commands that r in this file are not exposed to the user and can't be called

import (
	"fmt"

	"github.com/seme4eg/sadbot/stream"
)

// reply to context channel with message msg. Delete previous bot message to
// not flood the channel.
// NOTE: adding an option to not delete previous message will currently break
// pagination package since the latter stops observing reaction events on
// message delete event.
func (c *Ctx) reply(msg string) {

	// retrieve 10 previous messages in given channel
	messages, err := c.S.ChannelMessages(c.M.ChannelID, 10, "", "", "")
	if err != nil {
		fmt.Println("Error retrieving channel messages:", err)
	}

	// delete last bot message (if found) to not flood the channel
	// first message in slice is last in channel
	for _, m := range messages {
		if m.Author.ID == c.S.State.User.ID {
			err := c.S.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				fmt.Println("Error deleting previous message:", err)
			}
			break
		}
	}

	_, err = c.S.ChannelMessageSend(c.M.ChannelID, msg)
	if err != nil {
		fmt.Println("Failed to send message to channel:", err)
	}
}

// join joins voice, sets ctx VoiceConnection to streams map.
// For now is called only by Play and Playfolder f-s.
func join(ctx Ctx) error {
	// Get the voice state for the given guild and user
	VoiceState, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	c, err := ctx.S.ChannelVoiceJoin(ctx.M.GuildID, VoiceState.ChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice: ", err)
		return err
	}

	// create new stream and add it to pool of all streams with given guild id
	ctx.Streams[ctx.M.GuildID] = stream.New(c)

	return nil
}

// requirePresence requires presence of user and bot in same voice channel.
// Requires user and bot to be in the SAME channel.
// Returns error if any condition isn't met.
func requirePresence(ctx Ctx) error {
	// Get the voice state for the given guild and user
	VoiceState, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)

	// if err means user is not connected to a voice channel
	if err != nil {
		ctx.reply("Must be connected to voice channel to use bot")
		return err
	}

	// if bot is not connected to a voice channel in current guild..
	if ctx.stream().V == nil {
		ctx.reply(fmt.Sprintf(
			"Bot is not connected to a voice channel. Use **%splay <Title | URL>**",
			ctx.Prefix))
		return fmt.Errorf("bot is not connected to voice channel")
	}

	// forbid user to command bot if in different channel than the bot
	if ctx.stream().V.ChannelID != VoiceState.ChannelID {
		ctx.reply("Must be in same voice channel as bot")
		return fmt.Errorf("user is in different channel than bot")
	}

	return nil
}
