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
	// initial page paginator will start with, is determined by current track index
	var page int

	for i := 0; i < len(ctx.stream().Queue); i += perPage {

		embed := &discordgo.MessageEmbed{Title: "Current queue"}

		// maximum amount of tracks on current paginator page
		maxLen := int(math.Min(
			float64(i+perPage),
			float64(len(ctx.stream().Queue)-1)))

		for j := i; j <= maxLen; j++ {
			if ctx.stream().Queue[j].IsPlaying() {
				page = j / perPage
			}

			field := getField(j, ctx)
			embed.Fields = append(embed.Fields, field)
		}

		// to start showing next page from new index and not repeat last track
		// from previous page
		i++

		p.Add(embed)
	}

	// REVIEW: should query message also be some kind of a 'player controller' ?
	// i can add custom handlers here like play / pause / stop / clear / save
	// p.Handle("+", func() { ctx.S.ChannelMessageSend(ctx.M.ChannelID, "Hi!") })

	p.Spawn(page)
}

func getField(index int, ctx Ctx) *discordgo.MessageEmbedField {
	// needed for further format function, whether to number tracks 01 or 001
	var (
		numLen         = len(fmt.Sprint(len(ctx.stream().Queue)))
		formatNum      = "%0" + fmt.Sprint(numLen) + "d "
		formattedIndex = fmt.Sprintf(formatNum, index+1)
	)
	return &discordgo.MessageEmbedField{
		Value: fmt.Sprintf("%s – %s", formattedIndex, ctx.stream().Queue[index]),
	}
}
