package stream

import (
	"fmt"
	"math/rand"
	"sadbot/utils"

	"github.com/bwmarrin/discordgo"
)

type RepeatState string

const (
	RepeatOff    RepeatState = "off"
	RepeatSingle RepeatState = "single"
	RepeatAll    RepeatState = "all"
)

// didn't go for interface here cuz it seemed redundant
// 1 stream per 1 session, holds info about voice connection state, and other
// 'player' stuff like queue etc
type Stream struct {
	S *discordgo.Session
	V *discordgo.VoiceConnection

	Queue     []Song
	SongIndex uint16

	Playing bool
	Repeat  RepeatState

	Shuffled      bool
	ShuffledQueue []Song

	Stop chan bool
}

type Song struct {
	Title  string
	Source string
	// Duration string
}

func (s *Stream) Play() error {
	s.Playing = true

	for len(s.Queue) > 0 {
		s.SongIndex = 0
		if s.Shuffled {
			s.SongIndex = uint16(rand.Intn(len(s.Queue)))
		}

		song := s.Queue[s.SongIndex]

		fmt.Println("PlayAudioFile:", song.Title)
		s.S.UpdateWatchStatus(0, song.Title)

		err := utils.PlayAudioFile(s.V, song.Source, s.Stop, &s.Playing)
		if err != nil {
			fmt.Println("Error playing audio file: ", err)
			s.Playing = false
			return err
		}

		if len(s.Queue) > 0 {
			s.Queue = s.Queue[1:]
		}

		// 'next' command doesn't set this flag to 'false' but 'stop' does so we
		// need to check it here
		if !s.Playing {
			break
		}
	}

	s.Playing = false

	return nil
}

// TODO: on stop, leave
func (s *Stream) Reset() {
	s.Queue = s.Queue[:0]
	s.Playing = false
	s.Stop <- true
}

func (s *Stream) Shuffle() {
	s.Shuffled = true
	copy(s.ShuffledQueue, s.Queue)
	rand.Shuffle(len(s.ShuffledQueue), func(i, j int) {
		s.ShuffledQueue[i], s.ShuffledQueue[j] = s.ShuffledQueue[j], s.ShuffledQueue[i]
	})
}

func (s *Stream) UnShuffle() {
	s.Shuffled = false
	s.ShuffledQueue = s.ShuffledQueue[:0]
}

func (s *Stream) SetRepeat(state string) (response string) {
	switch RepeatState(state) {
	case RepeatSingle:
		s.Repeat = RepeatSingle
		return "Now repeating " + s.Queue[s.SongIndex].Title
	case RepeatAll:
		s.Repeat = RepeatAll
		// no matter len of what we take here - shuffled queue or not, length's same
		return "Now repeating " + fmt.Sprint(len(s.Queue)) + " songs"
	case RepeatOff:
		s.Repeat = RepeatOff
		return "Repeating turned off"
	default:
		return "Usage: <prefix>repeat single | all | off"
	}

}
