package cmds

import "fmt"

var cmDocs = []struct{ Command, Doc string }{
	{"help", "Show this help message."},
	{"Ping", "Ping the bot, checking if it's up and running."},
	{"p | play [<Title | URL>]", "Resume paused playback. If <title | url> argument was given plays links from youtube & soundcloud."},
	{"pf | playfolder <path>", "Play music from local folder (bot hoster only)."},
	{"pause", "Pause current playback."},
	{"stop", "Stop current playback, clears queue."},
	{"next | skip", "Skip to next song assuming there is one."},
	{"skipto <index>", "Skip to song with <index> (look it up in queue) assuming there is one."},
	{"clear", "Clear current queue."},
	{"leave", "Leave the channel, clears queue, resets shuffle & repeat states."},
	{"repeat single/all/off", "Repeat current track. Repeat all queue. Disable repeat."},
	{"loop", "Same as _repeat all_."},
	{"(un)shuffle", "(un)Shuffle the queue"},
	{"q | queue", "Show current queue."},
	{"np", "~~No problem~~ Show current track name."},
}

func Help(ctx Ctx) {
	var msg string
	for _, v := range cmDocs {
		msg += fmt.Sprintf("**%s%s** — %s\n", ctx.Prefix, v.Command, v.Doc)
	}

	// add footer note
	msg += "\nBot invite link:\n" + "https://discord.com/api/oauth2/authorize?client_id=1104687184537190441&permissions=274881440832&scope=bot"

	ctx.Reply(msg)
}
