package cmd

import (
	"fmt"
	"github.com/seme4eg/sadbot/utils"
	"strings"
)

// Play joins bot to voice if it is not in one. If no arguments passed calls
// current stream 'resume' method. Otherwise processes given query, calls
// current stream Add method for each processed track. When done replies with
// currnet queue. Then calls for Play method.
func Play(ctx Ctx) {
	// Get the voice state for the given guild and user
	_, err := ctx.S.State.VoiceState(ctx.M.GuildID, ctx.M.Author.ID)

	// if err means user is not connected to a voice channel
	if err != nil {
		ctx.reply("Must be connected to voice channel to use bot")
		return
	}

	// join voice in case bot is not in one
	if ctx.stream().V == nil {
		err := join(ctx)
		if err != nil {
			fmt.Println("Failed to join voice channel:", err)
			return
		}
	}

	args := strings.TrimSpace(ctx.Args)

	if args == "" {
		ctx.stream().Resume()
		return
	}

	res, err := utils.ProcessQuery(args)
	if err != nil {
		ctx.reply(err.Error())
	}

	for _, t := range res {
		ctx.stream().Add(t.Url, t.Title)
	}

	go Queue(ctx)

	if err := ctx.stream().Play(); err != nil {
		fmt.Println("Error streaming:", err)
	}
}
