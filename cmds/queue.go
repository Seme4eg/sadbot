package cmds

import "strings"

func Queue(ctx Ctx) {
	// REVIEW: should this command be accessible to only those in voice chat or no
	// err := RequirePresence(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// TODO: beautify that thing
	// - add nice title "Current queue:"
	// make it a complex send

	if len(ctx.Stream.Queue) == 0 {
		ctx.Reply("Queue is empty, sir")
	} else {
		trackNames := make([]string, len(ctx.Stream.Queue))
		for i, v := range ctx.Stream.Queue {
			trackNames[i] = v.Title
		}
		ctx.Reply(strings.Join(trackNames, "\n"))
	}
}
