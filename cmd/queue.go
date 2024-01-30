package cmd

import (
	"fmt"
	"github.com/seme4eg/sadbot/paginator"
	"math"

	"github.com/bwmarrin/discordgo"
)

// Queue replies to user with current stream queue. Reply message is paginator.
// Formats currently playing track with bold.
func Queue(ctx Ctx) {
	if len(ctx.stream().Queue) == 0 {
		ctx.reply("Queue is empty, sir")
		return
	}

	p := paginator.NewPaginator(ctx.S, ctx.M.ChannelID)

	// 10 tracks per page (starting with 0 index)
	perPage := 9
	// needed for further format function, whether to number tracks 01 or 001
	numLen := len(fmt.Sprint(len(ctx.stream().Queue)))
	formatNum := "%0" + fmt.Sprint(numLen) + "d "
	// initial page paginator will start with, is determined by current song index
	var page int

	for i := 0; i < len(ctx.stream().Queue); i += perPage {

		embed := &discordgo.MessageEmbed{Title: "Current queue"}

		// maximum amount of tracks on current paginator page
		maxLen := int(math.Min(
			float64(i+perPage),
			float64(len(ctx.stream().Queue)-1)))

		for j := i; j <= maxLen; j++ {
			formatName := "**%s** – %s"
			// change page & format for current song
			if j == ctx.stream().SongIndex {
				formatName = "**%s** – **%s**"
				page = j / perPage
			}
			field := discordgo.MessageEmbedField{
				Value: fmt.Sprintf(formatName,
					fmt.Sprintf(formatNum, j+1), ctx.stream().Queue[j].Title),
			}
			embed.Fields = append(embed.Fields, &field)
		}

		// to start showing next page from new index and not repeat last song
		// from previous page
		i++

		p.Add(embed)
	}

	// REVIEW: should query message also be some kind of a 'player controller' ?
	// i can add custom handlers here like play / pause / stop / clear / save
	// p.Handle("+", func() { ctx.S.ChannelMessageSend(ctx.M.ChannelID, "Hi!") })

	p.Spawn(page)
}
