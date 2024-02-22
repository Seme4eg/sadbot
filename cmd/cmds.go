package cmd

import (
	"strings"

	"github.com/seme4eg/sadbot/stream"

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
	Streams stream.Streams
	Prefix  string
}

// return stream of the guild in which message event happened
func (c *Ctx) stream() *stream.Stream {
	if stream, ok := c.Streams[c.M.GuildID]; ok {
		return stream
	}
	// if no stream in given guild - return new one just to have access to methods
	// like queue
	return &stream.Stream{}
}

// Handle handles user commands received with Discord 'message create' event.
// Creates new Context struct and passes it to handling function (if one found)
func Handle(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	streams stream.Streams,
	prefix string,
) {
	command := strings.TrimPrefix(m.Content, prefix)
	if len(command) < 1 {
		return
	}
	name := strings.Fields(command)[0]
	args := strings.TrimPrefix(command, name)
	args = strings.TrimSpace(args)

	if command, ok := commands[strings.ToLower(name)]; ok {
		ctx := Ctx{S: s, M: m, Args: args, Streams: streams, Prefix: prefix}
		command(ctx)
	}
}
