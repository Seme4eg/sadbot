package cmds

import (
	"fmt"
	_ "io/ioutil"
	_ "sadbot/utils"
)

func Play(ctx Ctx) {
	if len(ctx.Args) == 0 {
		ctx.Reply("link(s) pls, sir") // TODO: (s)
		return
	}

	err := Join(ctx)
	if err != nil {
		fmt.Println("Failed to join voice channel:", err)
		return
	}

	ctx.Stream.Playing = true

	// ctx.Reply("Adding songs to queue...") // XXX: maybe remove it

	// Start loop and attempt to play all files in the given folder
	// fmt.Println("Reading Folder: ", ctx.Args[0])
	// files, _ := ioutil.ReadDir(ctx.Args[0])
	// for _, f := range files {
	// 	fmt.Println("PlayAudioFile:", f.Name())
	// 	ctx.S.UpdateWatchStatus(0, f.Name())

	// 	err := utils.PlayAudioFile(vc,
	// 		fmt.Sprintf("%s/%s", ctx.Args[0], f.Name()), make(chan bool))
	// 	if err != nil {
	// 		fmt.Println("Error playing audio file: ", err)
	// 	}
	// }

	// for _, arg := range ctx.Args {
	// 	t, inp, err := ctx.Youtube.Get(arg)

	// 	if err != nil {
	// 		ctx.Reply("Rip, error, c logs.")
	// 		fmt.Println("Error getting input:", err)
	// 		return
	// 	}

	// 	switch t {
	// 	case framework.ERROR_TYPE:
	// 		ctx.Reply("An error occured!")
	// 		fmt.Println("error type", t)
	// 		return
	// 	case framework.VIDEO_TYPE:
	// 		{
	// 			video, err := ctx.Youtube.Video(*inp)
	// 			if err != nil {
	// 				ctx.Reply("An error occured!")
	// 				fmt.Println("error getting video1,", err)
	// 				return
	// 			}
	// 			song := framework.NewSong(video.Media, video.Title, arg)
	// 			sess.Queue.Add(*song)
	// 			ctx.Discord.ChannelMessageEdit(ctx.TextChannel.ID, msg.ID, "Added `"+song.Title+"` to the song queue."+
	// 				" Use `music play` to start playing the songs! To see the song queue, use `music queue`.")
	// 			break
	// 		}
	// 	case framework.PLAYLIST_TYPE:
	// 		{
	// 			videos, err := ctx.Youtube.Playlist(*inp)
	// 			if err != nil {
	// 				ctx.Reply("An error occured!")
	// 				fmt.Println("error getting playlist,", err)
	// 				return
	// 			}
	// 			for _, v := range *videos {
	// 				id := v.Id
	// 				_, i, err := ctx.Youtube.Get(id)
	// 				if err != nil {
	// 					ctx.Reply("An error occured!")
	// 					fmt.Println("error getting video2,", err)
	// 					continue
	// 				}
	// 				video, err := ctx.Youtube.Video(*i)
	// 				if err != nil {
	// 					ctx.Reply("An error occured!")
	// 					fmt.Println("error getting video3,", err)
	// 					return
	// 				}
	// 				song := framework.NewSong(video.Media, video.Title, arg)
	// 				sess.Queue.Add(*song)
	// 			}
	// 			ctx.Reply("Finished adding songs to the playlist. Use `music play` to start playing the songs! " +
	// 				"To see the song queue, use `music queue`.")
	// 			break
	// 		}
	// 	}
	// }
}
