package cmd

import "fmt"

// Pause calls current guild's stream Pause method.
func Pause(ctx Ctx) {
	err := requirePresence(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.stream().Pause()
}
