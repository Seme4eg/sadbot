package cmds

import (
	"fmt"
	"strconv"
)

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

	// -1 so in stream method we work with 0-based indecies
	if err := ctx.Stream.Skipto(index - 1); err != nil {
		ctx.Reply(err.Error())
	} else {
		ctx.Reply("Now playing: " + ctx.Stream.Current())
	}
}
