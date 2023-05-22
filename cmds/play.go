package cmds

import (
	"fmt"
	"github.com/seme4eg/sadbot/utils"
	"strings"
)

// Play joins bot to voice if it is not in one. If no arguments passed calls
// current stream Unpause method. Otherwise processes given query, calls current
// stream Add method for each processed track. When done replies with currnet
// queue. Then calls for Play method.
func Play(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// join voice in case bot is not in one
	if ctx.Stream().V == nil {
		err := Join(ctx)
		if err != nil {
			fmt.Println("Failed to join voice channel:", err)
			return
		}
	}

	args := strings.TrimSpace(ctx.Args)

	if args == "" {
		ctx.Stream().Unpause()
		return
	}

	res, err := utils.ProcessQuery(args)
	if err != nil {
		ctx.Reply(err.Error())
	}

	for _, t := range res {
		ctx.Stream().Add(t.Url, t.Title)
	}

	go Queue(ctx)

	if err := ctx.Stream().Play(); err != nil {
		fmt.Println("Error streaming:", err)
	}
}
