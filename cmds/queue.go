package cmds

func Queue(ctx Ctx) {
	// REVIEW: should this command be accessible to only those in voice chat or no
	// err := RequirePresence(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	if len(ctx.Stream.Queue) == 0 {
		ctx.Reply("Queue is empty, sir")
	}
}
