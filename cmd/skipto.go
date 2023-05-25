package cmd

import (
	"fmt"
	"strconv"
)

// Skipto checks if passed index is valid (convertable to int). If not
// replies with correct command usage. Otherwise calls for current guild's
// stream Skipto method outputting current song on success or error otherwise.
func Skipto(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	index, err := strconv.Atoi(ctx.Args)
	if err != nil {
		ctx.Reply("Usage: **skipto 5** where '5' is a song index that you can find with `queue` command.")
		return
	}

	// -1 to make index 0-based
	if err := ctx.Stream().Skipto(index - 1); err != nil {
		ctx.Reply(err.Error())
	} else {
		ctx.Reply("Now playing: " + ctx.Stream().Current())
	}
}
