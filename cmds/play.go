package cmds

import (
	"fmt"
	"sadbot/utils"
	"strings"
)

func Play(ctx Ctx) {
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

	args := strings.TrimSpace(ctx.Args)

	if args == "" {
		ctx.Stream.Playing = true
		return
	}

	res, err := utils.ProcessQuery(args)
	if err != nil {
		ctx.Reply(err.Error())
	}

	for _, t := range res {
		ctx.Stream.Add(t.Url, t.Title)
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
