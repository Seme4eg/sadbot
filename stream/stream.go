package stream

import (
	"fmt"
	"math/rand"
	"sadbot/utils"
	"sort"

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

	// NOTE: when adding new field don't forget to reset it on events like
	// leave, stop and clear

	S *discordgo.Session
	V *discordgo.VoiceConnection

	Queue     []Song
	SongIndex int

	Playing bool
	Repeat  RepeatState

	// after shuffle / unshuffle commands queue changes and current song index
	// in queue is already different song, so in order to not skip over it
	// this flag is needed
	SkipIndexUpdate bool

	Stop chan bool
}

type Song struct {
	Title  string
	Source string
	Index  int // initial index, used to 'unshuffle' suffled queue
	// TODO: Duration string
}

func (s *Stream) Play() error {
	s.Playing = true

	for len(s.Queue) > 0 {
		s.SongIndex = 0

		song := s.Queue[s.SongIndex]

		fmt.Println("PlayAudioFile:", song.Title)
		s.S.UpdateListeningStatus(song.Title)

		err := utils.PlayAudioFile(s.V, song.Source, s.Stop, &s.Playing)
		if err != nil {
			fmt.Println("Error playing audio file: ", err)
			s.Playing = false
			return err
		}

		if s.Repeat == RepeatOff && len(s.Queue) > 0 {
			s.Queue = s.Queue[1:]
		} else if s.Repeat == RepeatAll {
			if s.SongIndex+1 >= len(s.Queue) {
				s.SongIndex = 0
			} else {
				s.SongIndex++
			}
		}
		// in case of RepeatSingle nothing changes and song index remains same

		// 'next' command doesn't set this flag to 'false' but 'stop' does so we
		// need to check it here
		if !s.Playing {
			break
		}
	}

	s.Playing = false

	return nil
}

// Resets all fields of current (and only one) Stream struct except
// field 'S *discordgo.Session' and Stop channel
// and stops possibly remaining ffmpeg process
// withoutVoiceChan flag is passed for now only when using Stop command
// since it needs to reset Stream in the same way except Voice Channel
// cuz it's still in one
func (s *Stream) Reset(withoutVoiceChan bool) {
	s.Stop <- true
	if !withoutVoiceChan {
		s.V = nil
	}
	s.Queue = []Song{}
	s.SongIndex = 0
	s.Playing = false
	s.Repeat = RepeatOff
	s.SkipIndexUpdate = false
}

// randomise queue except 1st song
func (s *Stream) Shuffle() {
	currentSong := s.Queue[s.SongIndex]
	// create temporary 'queue' value that doesn't contain currently playing
	// track since after each shuffle we want it to be still first in queue
	temp := append(s.Queue[:s.SongIndex], s.Queue[s.SongIndex+1:]...)
	rand.Shuffle(len(temp), func(i, j int) {
		temp[i], temp[j] = temp[j], temp[i]
	})
	s.Queue = append([]Song{currentSong}, temp...)
}

func (s *Stream) UnShuffle() {
	// sort queue in ascending order by index field
	sort.Slice(s.Queue, func(i, j int) bool {
		return s.Queue[i].Index < s.Queue[j].Index
	})
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

func (s *Stream) Next() string {
	if s.SongIndex+1 >= len(s.Queue) {
		return "Either last song in the queue or no songs in it"
	}

	s.Playing = true
	// sends stop signal to current ffmpeg command stopping it
	s.Stop <- true
	return ""
}
