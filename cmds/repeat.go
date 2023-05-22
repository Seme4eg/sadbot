package cmds

import "fmt"

// Repeat sets current guild's stream repeat state to either single / all or off.
// In case of invalid argument passed replies with correct usage message.
// On repeat 'all' event replies with the amount of songs in current repeat loop.
// On repaet 'single' event replies with the name of the track that is on repeat.
func Repeat(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ctx.Stream().SetRepeat(ctx.Args)
	if err != nil {
		ctx.Reply(fmt.Sprintf("Usage: **%srepeat single | all | off**", ctx.Prefix))
		return
	}
	// kinda hardcode but user responses should be handled here, not in Stream
	switch ctx.Args {
	case "single":
		ctx.Reply(fmt.Sprintf("Now repeating: **%s**", ctx.Stream().Current()))
	case "all":
		ctx.Reply(
			fmt.Sprintf("Now repeating **%s** songs", fmt.Sprint(len(ctx.Stream().Queue))))
	case "off":
		ctx.Reply("Repeat turned off")
	}
}
