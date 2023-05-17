package stream

import (
	"fmt"
	"math/rand"
	"sadbot/utils"
	"sort"
	"sync"

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
	sync.Mutex
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

func New() *Stream {
	return &Stream{
		Stop:   make(chan bool),
		Repeat: RepeatOff,
	}
}

func (s *Stream) Play() error {
	s.Playing = true

	for len(s.Queue) > 0 && s.SongIndex >= 0 && s.SongIndex < len(s.Queue) && s.Playing {

		song := s.Queue[s.SongIndex]

		done := make(chan bool)
		defer close(done)
		go func() {
			err := utils.PlayAudioFile(s.V, song.Source, s.Stop, &s.Playing)
			if err != nil {
				fmt.Println("Error playing audio file: ", err)
			}
			close(done)
		}()

		select {
		// in case play function finished playing on its own (wasn't affected by
		// user commands) - skip to next song
		case <-done:
			fmt.Println("entered done case")
			if s.Repeat != RepeatSingle {
				if err := s.Next(); err != nil {
					fmt.Println("error nexting:", err)
					s.Playing = false
					return err
				}
			}
		case <-s.Stop:
			s.Stop <- true // also stop current ffmpeg process
			// prevent ffmpeg processes from overlapping (especially on prev command)
			time.Sleep(time.Millisecond * 350)
			// otherwise do nothing cuz it means user is skipping / preving
			continue
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
	s.Lock()
	defer s.Unlock()
	currentSong := s.Queue[s.SongIndex]
	// create temporary 'queue' value that doesn't contain currently playing
	// track since after each shuffle we want it to be still first in queue
	temp := append(s.Queue[:s.SongIndex], s.Queue[s.SongIndex+1:]...)
	rand.Shuffle(len(temp), func(i, j int) {
		temp[i], temp[j] = temp[j], temp[i]
	})
	s.Queue = append([]Song{currentSong}, temp...)
}

// sorts queue based on songs initial index
func (s *Stream) UnShuffle() {
	s.Lock()
	defer s.Unlock()
	// sort queue in ascending order by index field
	sort.Slice(s.Queue, func(i, j int) bool {
		return s.Queue[i].Index < s.Queue[j].Index
	})
}

// sets stream repeat state and returns response string
func (s *Stream) SetRepeat(state string) (response string) {
	s.Lock()
	defer s.Unlock()
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

func (s *Stream) Next() error {
	s.Lock()

	s.Playing = false
	s.SongIndex++

	if s.SongIndex >= len(s.Queue) {
		switch s.Repeat {
		case RepeatAll:
			s.SongIndex = 0
		case RepeatOff:
			s.SongIndex--
			return errors.New("either last song in the queue or no songs in it")
		}
	}

	s.Playing = true
	s.Unlock()
	// sends stop signal to current ffmpeg command stopping it
	s.Stop <- true
	return nil
}

func (s *Stream) Prev() error {
	s.Lock()

	s.Playing = false
	s.SongIndex--

	if s.SongIndex < 0 {
		switch s.Repeat {
		case RepeatAll:
			s.SongIndex = len(s.Queue) - 1
		case RepeatOff:
			s.SongIndex++
			return errors.New("nothing was played before")
		}
	}

	s.Playing = true
	s.Unlock()
	// sends stop signal to current ffmpeg command stopping it
	s.Stop <- true
	return nil
}

func (s *Stream) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Queue = s.Queue[:0]
	s.SongIndex = 0
}

func (s *Stream) Current() string {
	if len(s.Queue) == 0 {
		return ""
	}
	return s.Queue[s.SongIndex].Title
}

func (s *Stream) Add(Source, Title string) {
	s.Lock()
	defer s.Unlock()
	s.Queue = append(s.Queue, Song{Title, Source, len(s.Queue)})
}
