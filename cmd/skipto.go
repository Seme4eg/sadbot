package cmd

import (
	"fmt"
	"strconv"
)

// Skipto checks if passed index is valid (convertable to int). If not
// replies with correct command usage. Otherwise calls for current guild's
// stream Skipto method outputting current track on success or error otherwise.
func Skipto(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	index, err := strconv.Atoi(ctx.Args)
	if err != nil {
		ctx.reply("Usage: **skipto 5** where '5' is a track index that you can find with `queue` command.")
		return
	}

	// -1 to make index 0-based
	if err := ctx.stream().Skipto(index - 1); err != nil {
		ctx.reply(err.Error())
	} else {
		ctx.reply("Now playing: " + ctx.stream().Current())
	}
}
