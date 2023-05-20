package cmds

import "fmt"

func Repeat(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ctx.Stream.SetRepeat(ctx.Args)
	if err != nil {
		ctx.Reply(fmt.Sprintf("Usage: **%srepeat single | all | off**", ctx.Prefix))
		return
	}
	// kinda hardcode but user responses should be handled here, not in Stream
	switch ctx.Args {
	case "single":
		ctx.Reply(fmt.Sprintf("Now repeating: **%s**", ctx.Stream.Current()))
	case "all":
		ctx.Reply(
			fmt.Sprintf("Now repeating **%s** songs", fmt.Sprint(len(ctx.Stream.Queue))))
	case "off":
		ctx.Reply("Repeat turned off")
	}
}
