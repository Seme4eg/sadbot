package cmds

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sadbot/stream"
	"strings"
)

// plays files from local folder (for now doesn't support standalones for now)
func PlayFolder(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// join voice in case bot is not in one
	if ctx.Stream.V == nil {
		err := Join(ctx)
		if err != nil {
			fmt.Println("Failed to join voice channel:", err)
			return
		}
	}

	if strings.TrimSpace(ctx.Args) == "" {
		ctx.Reply("provide a folder pls, sir")
		return
	}

	ctx.Reply("Adding songs to queue...") // XXX: maybe remove it

	// Start loop and attempt to play all files in the given folder
	fmt.Println("Reading Folder: ", ctx.Args)

	// TODO: output here which tracks got added to queue

	files, err := ioutil.ReadDir(filepath.FromSlash(ctx.Args))
	if err != nil {
		fmt.Println("ReadDir error:", err)
		return
	}

	// add tracks to queue
	for _, f := range files {
		// TODO: not sure is Song type must b part of 'stream' module
		ctx.Stream.Queue = append(ctx.Stream.Queue, stream.Song{
			Title: f.Name(),
			// source - full path to the file
			Source: ctx.Args + "/" + f.Name(),
		})
	}

	err = ctx.Stream.Play()
	if err != nil {
		fmt.Println("Error streaming:", err)
		return
	}

}
