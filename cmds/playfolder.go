package cmds

import (
	"fmt"
	"path/filepath"
	"strings"
)

var formats = []string{"mp3", "flac", "wav", "opus"}

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
		ctx.Stream.Add(path, title)
	}

	go Queue(ctx)

	ctx.S.UpdateListeningStatus(ctx.Stream.Current())
	err = ctx.Stream.Play()
	if err != nil {
		fmt.Println("Error streaming:", err)
		return
	}
	// TODO: check if it updates to help
	ctx.S.UpdateListeningStatus(ctx.Prefix + "help")
}
