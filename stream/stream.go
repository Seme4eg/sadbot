package stream

import (
	"fmt"
	"sadbot/utils"

	"github.com/bwmarrin/discordgo"
)

// didn't go for interface here cuz it seemed redundant
// 1 stream per 1 session, holds info about voice connection state, and other
// 'player' stuff like queue etc
type Stream struct {
	S       *discordgo.Session
	V       *discordgo.VoiceConnection
	Queue   []Song
	Playing bool
	Stop    chan bool
}

// XXX: should it be private?
type Song struct {
	Title  string
	Source string
	// Duration string
}

func (s *Stream) Play() error {
	s.Playing = true

	fmt.Println(s.Queue)

	for len(s.Queue) > 0 {
		song := s.Queue[0]
		fmt.Println("PlayAudioFile:", song.Title)
		s.S.UpdateWatchStatus(0, song.Title)
		err := utils.PlayAudioFile(s.V, song.Source, s.Stop)
		if err != nil {
			fmt.Println("Error playing audio file: ", err)
			s.Playing = false
			return err
		}
		s.Queue = s.Queue[1:]
		if !s.Playing {
			break
		}
	}

	s.Playing = false

	return nil
}

func (s *Stream) Next() {
	fmt.Println("sending stop signal")
	s.Playing = true
	s.Stop <- true
}
