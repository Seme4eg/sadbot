package cmds

import "fmt"

// Pause calls current guild's stream Pause method.
func Pause(ctx Ctx) {
	err := RequirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Stream().Pause()
}
