package cmds

import (
	"fmt"
	"sadbot/stream"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commands = map[string]func(Ctx){
	"ping": func(ctx Ctx) {
		_, _ = ctx.S.ChannelMessageSend(ctx.M.ChannelID, "pong")
	},
	"play":       Play,
	"p":          Play,
	"playfolder": PlayFolder,
	"pf":         PlayFolder,
	"pause":      Pause,
	"stop":       Stop,
	"next":       Next,
	"skip":       Next,
	"skipto":     Skipto,
	"prev":       Prev,
	"clear":      Clear,
	"leave":      Leave,
	"repeat":     Repeat,
	"loop": func(ctx Ctx) {
		ctx.Args = "all"
		Repeat(ctx)
	},
	"shuffle":   Shuffle,
	"unshuffle": Unshuffle,
	"queue":     Queue,
	"q":         Queue,
	"np":        NowPlaying,
	"help":      Help,
}

// context for each separate message adressed to bot
type Ctx struct {
	S       *discordgo.Session
	M       *discordgo.MessageCreate
	Args    string
	Streams *stream.Streams
	Prefix  string
}

// XXX: maybe i don't need to return empty one? maybe there is a way to
// always preserve user from not having stream or smth else
// simple stream getter for the guild in which message event happened
func (c *Ctx) Stream() *stream.Stream {
	if stream, ok := c.Streams.List[c.M.GuildID]; ok {
		return stream
	}
	// if no stream in given guild - return new one just to have access to methods
	// like queue
	return &stream.Stream{}
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

// Commands that r in this file are not exposed to the user and can't be called

// joins voice, sets ctx VoiceConnection to streams map
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

	ctx.Streams.List[ctx.M.GuildID] = stream.New(c)

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

	if ctx.Stream().V != nil {
		// forbid user to command 'leave' if in different channel than the bot
		if ctx.Stream().V.ChannelID != VoiceState.ChannelID {
			ctx.Reply("Must be in same voice channel as bot")
			return fmt.Errorf("user is in different channel than bot")
		}
	}

	return nil
}

func Handle(command string, s *discordgo.Session,
	m *discordgo.MessageCreate, streams *stream.Streams, prefix string) {
	name := strings.Fields(command)[0]
	args := strings.TrimPrefix(command, name+" ")

	if command, ok := commands[strings.ToLower(name)]; ok {
		ctx := Ctx{S: s, M: m, Args: args, Streams: streams, Prefix: prefix}
		command(ctx)
	}
}
