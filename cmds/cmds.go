package cmds

import (
	"fmt"
	"github.com/seme4eg/sadbot/stream"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Map of aliases for bot commands.
var commands = map[string]func(Ctx){
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

// Ctx is a context for each separate user message
type Ctx struct {
	S       *discordgo.Session
	M       *discordgo.MessageCreate
	Args    string
	Streams *stream.Streams
	Prefix  string
}

// return Stream of the guild in which message event happened
func (c *Ctx) Stream() *stream.Stream {
	if stream, ok := c.Streams.List[c.M.GuildID]; ok {
		return stream
	}
	// if no stream in given guild - return new one just to have access to methods
	// like queue
	return &stream.Stream{}
}

// Reply to context channel with message msg. Delete previous bot message to
// not flood the channel.
// NOTE: adding an option to not delete previous message will currently break
// pagination package since the latter stops observing reaction events on
// message delete event.
func (c *Ctx) Reply(msg string) {

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

// Commands that r in this file are not exposed to the user and can't be called

// Join joins voice, sets ctx VoiceConnection to streams map.
// For now is called only by Play and Playfolder f-s.
func Join(ctx Ctx) error {
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
	ctx.Streams.List[ctx.M.GuildID] = stream.New(c)

	return nil
}

// RequirePresence requires presence of user and bot in same voice channel.
// Requires user and bot to be in the SAME channel.
// Also currently this f-n is not called under conditions when bot is not in
// voice channel.
// Returns error if this condition isn't met.
func RequirePresence(ctx Ctx) error {
	// Get the voice state for the given guild and user
	VoiceState, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)

	// if err means user is not connected to a voice channel
	if err != nil {
		ctx.Reply("Must be connected to voice channel to use bot")
		return err
	}

	// if bot is not connected to a voice channel in current guild..
	if ctx.Stream().V == nil {
		ctx.Reply(fmt.Sprintf(
			"Bot is not connected to a voice channel. Use **%splay <Title | URL>**",
			ctx.Prefix))
		return fmt.Errorf("bot is not connected to voice channel")
	}

	// forbid user to command bot if in different channel than the bot
	if ctx.Stream().V.ChannelID != VoiceState.ChannelID {
		ctx.Reply("Must be in same voice channel as bot")
		return fmt.Errorf("user is in different channel than bot")
	}

	return nil
}

// Handle handles user commands received with Discord 'message create' event.
// Creates new Context struct and passes it to handling function (if one found)
func Handle(command string, s *discordgo.Session,
	m *discordgo.MessageCreate, streams *stream.Streams, prefix string) {
	name := strings.Fields(command)[0]
	args := strings.TrimPrefix(command, name+" ")

	if command, ok := commands[strings.ToLower(name)]; ok {
		ctx := Ctx{S: s, M: m, Args: args, Streams: streams, Prefix: prefix}
		command(ctx)
	}
}
