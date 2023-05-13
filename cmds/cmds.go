package cmds

import (
	"fmt"
	"sadbot/stream"

	"github.com/bwmarrin/discordgo"
)

// context for each separate message adressed to bot
type Ctx struct {
	S *discordgo.Session
	M *discordgo.MessageCreate
	// need args to be just a string rather than slice of strings cuz
	// there are functions (like playfolder) that don't need a pre-splitted args
	Args   string
	Stream *stream.Stream
}

func (c *Ctx) Reply(msg string) {
	messages, err := c.S.ChannelMessages(c.M.ChannelID, 10, "", "", "")
	if err != nil {
		fmt.Println("Error retrieving channel messages:", err)
	}

	// delete last bot message to not flood the channel
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

// most likely will change in future, object with aliases and commands
// it will evolve / change when i will b writing functionality for 'help' f-n
var Pool = map[string]func(ctx Ctx){
	"ping": func(ctx Ctx) {
		_, _ = ctx.S.ChannelMessageSend(ctx.M.ChannelID, "pong")
	},
	"leave":      Leave,
	"play":       Play,
	"p":          Play,
	"pause":      Pause,
	"playfolder": PlayFolder,
	"pf":         PlayFolder,
	"next":       Next,
	"skip":       Next,
	"clear":      Clear,
	"stop":       Stop,
	"repeat":     Repeat,
	"shuffle":    Shuffle,
	"queue":      Queue,
	"nowplaying": NowPlaying,
	"np":         NowPlaying,
}

// Commands that r in this file are not exposed to the user and can't be called

// joins voice, sets ctx.Stream.V(oiceConnection)
func Join(ctx Ctx) error {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

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
	ctx.Stream.V = c
	return nil
}

// require presence of user and bot in the SAME channel
// return error if this condition isn't met
func RequirePresence(ctx Ctx) error {
	// Get the voice state for the given guild and user
	VoiceState, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)

	// if err means user is not connected to a voice channel
	if err != nil {
		ctx.Reply("Must be connected to voice channel to use bot")
		return err
	}

	if ctx.Stream.V != nil {
		// forbid user to command 'leave' if in different channel than the bot
		if ctx.Stream.V.ChannelID != VoiceState.ChannelID {
			ctx.Reply("Must be in same voice channel as bot")
			return fmt.Errorf("user is in different channel than bot")
		}
	}

	return nil
}
