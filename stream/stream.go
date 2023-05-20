package stream

import (
	"errors"
	"fmt"
	"math/rand"
	"sadbot/utils"
	"sort"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type RepeatState string

const (
	RepeatOff    RepeatState = "off"
	RepeatSingle RepeatState = "single"
	RepeatAll    RepeatState = "all"
)

type Streams struct {
	// keys are guildid since it seems reasonable to store 1 stream per guild
	List map[string]*Stream
}

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

	Stop chan bool
}

type Song struct {
	Title  string
	Source string
	Index  int // initial index, used to 'unshuffle' suffled queue
	// TODO: Duration string
}

func New(vc *discordgo.VoiceConnection) *Stream {
	return &Stream{
		V:      vc,
		Stop:   make(chan bool),
		Repeat: RepeatOff,
	}
}

func (s *Stream) Play() error {
	s.Playing = true

	for len(s.Queue) > 0 && s.SongIndex >= 0 && s.SongIndex < len(s.Queue) && s.Playing {

		done := make(chan bool)
		// defer close(done)
		go func() {
			err := utils.PlayAudioFile(s.V, s.Queue[s.SongIndex].Source, s.Stop, &s.Playing)
			if err != nil {
				fmt.Println("Error playing audio file: ", err)
			}
			done <- true
		}()

		select {
		// in case play function finished playing on its own (wasn't affected by
		// user commands) - skip to next song
		case <-done:
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
func (s *Stream) Reset() {
	s.Lock()
	defer s.Unlock()
	s.Playing = false
	s.Queue = s.Queue[:0]
	s.SongIndex = 0
	s.Repeat = RepeatOff
	select {
	case s.Stop <- true:
	default:
		fmt.Println("Stop Channel is closed")
	}
}

func (s *Stream) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Queue = s.Queue[:0]
	s.SongIndex = 0
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
	s.SongIndex = 0
	fmt.Println(len(s.Queue), "len")
}

// sorts queue based on songs initial index
func (s *Stream) UnShuffle() {
	s.Lock()
	defer s.Unlock()
	// get index of currently playing song
	curIndex := s.Queue[s.SongIndex].Index

	// sort queue in ascending order by index field
	sort.Slice(s.Queue, func(i, j int) bool {
		return s.Queue[i].Index < s.Queue[j].Index
	})

	// upate currently playing song index
	for i, t := range s.Queue {
		if t.Index == curIndex {
			s.SongIndex = i
			break
		}
	}
}

// sets stream repeat state and returns response string
func (s *Stream) SetRepeat(state string) error {
	s.Lock()
	defer s.Unlock()
	switch val := RepeatState(state); val {
	case RepeatSingle, RepeatAll, RepeatOff:
		s.Repeat = val
	default:
		return errors.New("invalid repeat state passed")
	}
	return nil
}

// kills current playback, skips to next track
func (s *Stream) Next() error {
	s.Lock()
	defer s.Unlock()

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

	s.Stop <- true
	return nil
}

// kills current playback, skips to previous track
func (s *Stream) Prev() error {
	s.Lock()
	defer s.Unlock()

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

	s.Stop <- true
	return nil
}

func (s *Stream) Current() string {
	if len(s.Queue) == 0 {
		return ""
	}
	return s.Queue[s.SongIndex].Title
}

// TODO: add easyer support for more fields
func (s *Stream) Add(Source, Title string) {
	s.Lock()
	defer s.Unlock()
	s.Queue = append(s.Queue, Song{Title, Source, len(s.Queue)})
}

func (s *Stream) Skipto(index int) error {
	if index <= 0 || index > len(s.Queue) {
		return errors.New("no song with such index")
	}

	s.Lock()

	s.SongIndex = index

	s.Unlock()
	s.Stop <- true
	return nil
}

func (s *Stream) Disconnect() error {
	s.Reset()
	return s.V.Disconnect()
}

func (s *Stream) Pause() {
	s.Playing = false
}
