package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

var formats = []string{"mp3", "flac", "wav", "opus"}

// PlayFolder joins bot to voice if it is not in one. Processes given
// folder, calls current stream Add method for each processed track (file).
// When done replies with currnet queue. Then calls for Play method.
// (doesn't support standalones for now)
func PlayFolder(ctx Ctx) {
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

	if strings.TrimSpace(ctx.Args) == "" {
		ctx.reply("provide a folder pls, sir")
		return
	}

	fmt.Println("Reading Folder:", ctx.Args)

	var trackPaths []string
	// add files of each listed format to slice
	for _, f := range formats {
		path := filepath.Join(ctx.Args, "*."+f)
		paths, err := filepath.Glob(path)
		if err != nil {
			fmt.Println("Error getting", f, "files from dir:", err)
		}
		trackPaths = append(trackPaths, paths...)
	}

	// add tracks to queue
	for _, path := range trackPaths {
		title := strings.TrimPrefix(path, ctx.Args+"/")
		ctx.stream().Add(path, title)
	}

	go Queue(ctx)

	if err := ctx.stream().Play(); err != nil {
		fmt.Println("Error streaming:", err)
		return
	}
}
