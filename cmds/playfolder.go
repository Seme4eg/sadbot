package cmds

import (
	"fmt"
	"path/filepath"
	"sadbot/stream"
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

	// Start loop and attempt to play all files in the given folder
	fmt.Println("Reading Folder: ", ctx.Args)

	var trackPaths []string
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
		name := strings.TrimPrefix(path, ctx.Args+"/")
		ctx.Stream.Queue = append(ctx.Stream.Queue, stream.Song{
			Title:  name,
			Source: path, // full path to the file
			Index:  len(ctx.Stream.Queue),
		})
	}

	Queue(ctx)

	err = ctx.Stream.Play()
	if err != nil {
		fmt.Println("Error streaming:", err)
		return
	}

}
