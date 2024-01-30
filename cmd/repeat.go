package cmd

import "fmt"

// Repeat sets current guild's stream repeat state to either single / all or off.
// In case of invalid argument passed replies with correct usage message.
// On repeat 'all' event replies with the amount of songs in current repeat loop.
// On repaet 'single' event replies with the name of the track that is on repeat.
func Repeat(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ctx.stream().SetRepeat(ctx.Args)
	if err != nil {
		ctx.reply(fmt.Sprintf("Usage: **%srepeat single | all | off**", ctx.Prefix))
		return
	}
	// kinda hardcode but user responses should be handled here, not in Stream
	switch ctx.Args {
	case "single":
		ctx.reply(fmt.Sprintf("Now repeating: **%s**", ctx.stream().Current()))
	case "all":
		ctx.reply(
			fmt.Sprintf("Now repeating **%s** songs", fmt.Sprint(len(ctx.stream().Queue))))
	case "off":
		ctx.reply("Repeat turned off")
	}
}
