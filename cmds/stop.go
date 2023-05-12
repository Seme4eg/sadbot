package cmds

import "fmt"

func Stop(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	Clear(ctx)
	ctx.Stream.Stop <- true
}
