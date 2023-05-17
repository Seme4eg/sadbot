package cmds

import (
	"fmt"
	"sadbot/paginator"

	"github.com/bwmarrin/discordgo"
)

func Queue(ctx Ctx) {
	if len(ctx.Stream.Queue) == 0 {
		ctx.Reply("Queue is empty, sir")
		return
	}

	p := paginator.NewPaginator(ctx.S, ctx.M.ChannelID)

	// 10 tracks per page
	perPage := 10
	// needed for further format function, whether to number tracks 01 or 001
	numLen := len(fmt.Sprint(len(ctx.Stream.Queue)))
	formatNum := "%0" + fmt.Sprint(numLen) + "d "

	for i := 0; i < len(ctx.Stream.Queue); i += perPage {

		var maxLen int

		embed := &discordgo.MessageEmbed{
			Title: "Current queue",
			// Color: 0x00ff00,
		}

		if a, b := i+perPage, len(ctx.Stream.Queue)-1; a > b {
			maxLen = b
		} else {
			maxLen = a
		}

		for j := i; j <= maxLen; j++ {
			formatName := "**%s** – %s"
			if j == ctx.Stream.SongIndex {
				formatName = "**%s** – **%s**"
			}
			field := discordgo.MessageEmbedField{
				Value: fmt.Sprintf(formatName,
					fmt.Sprintf(formatNum, j+1), ctx.Stream.Queue[j].Title),
			}
			embed.Fields = append(embed.Fields, &field)
		}

		p.Add(embed)
	}

	// REVIEW: should query message also be some kind of a 'player' ?
	// i can add custom handlers here like play / pause / stop / clear / save
	// p.Handle("+", func() { ctx.S.ChannelMessageSend(ctx.M.ChannelID, "Hi!") })

	p.Spawn()
}
