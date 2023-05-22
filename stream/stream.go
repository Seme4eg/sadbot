package stream

import (
	"errors"
	"fmt"
	"github.com/seme4eg/sadbot/utils"
	"math/rand"
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

// needed to wrap the 'List' it in struct to not lose it's pointer when passing
// it around cuz then can't delete inactive streams
type Streams struct {
	// keys are guild ids since it seems reasonable to store 1 stream per guild
	List map[string]*Stream
}

// Stream holds 1 stream per 1 guild, holds info about voice connection state,
// and other 'player' stuff like queue etc
type Stream struct {
	// NOTE: when adding new field don't forget to reset it on events like
	// leave, stop and clear
	sync.Mutex
	V *discordgo.VoiceConnection

	Queue []Song
	// index of currently playing song in queue (not song initial index that
	// presents in Song struct)
	SongIndex int

	Playing bool
	Repeat  RepeatState

	// Channel which when being sent to stops current ffmpeg playback
	Stop chan bool
}

type Song struct {
	Title string
	// source can be either Url or a path (if using playfolder command)
	Source string
	Index  int // initial index, used to 'unshuffle' suffled queue
	// TODO: Duration string
}

// New returns new stream struct
func New(vc *discordgo.VoiceConnection) *Stream {
	return &Stream{
		V:      vc,
		Stop:   make(chan bool),
		Repeat: RepeatOff,
	}
}

// Play starts playback of song with current SongIndex and attempts to skip
// to next song in queue on current song end or on send event to Stop channel.
func (s *Stream) Play() error {
	s.Playing = true

	for len(s.Queue) > 0 && s.SongIndex >= 0 && s.SongIndex < len(s.Queue) && s.Playing {

		done := make(chan bool)
		// start playback function in background
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
			// this case can happen only on user iteration, means we do not need
			// to respect current Repeat state
			s.Stop <- true // stop current ffmpeg process
			// prevent ffmpeg processes from overlapping (especially on prev command)
			time.Sleep(time.Millisecond * 350)
			continue
		}
	}

	s.Playing = false

	return nil
}

// Reset resets fields of current Stream except session and stop channel.
// Stops possibly remaining ffmpeg process.
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

// Clear empties queue and resets current song index to 0.
func (s *Stream) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Queue = s.Queue[:0]
	s.SongIndex = 0
}

// Shuffle shuffles queue. Currently playing song will be 1st always in shuffled
// queue.
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

// Unshuffle sorts songs in queue by their initial index field.
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

// SetRepeat sets current guild's stream repeat state to either
// single / all or off.
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

// Next skips to next song unconditionally. Means even if user has set repeat
// state to 'single' it will still skip to next song. Kills current playback by
// sending to Stop channel.
func (s *Stream) Next() error {
	s.Lock()
	defer s.Unlock()

	// FIXME: somewhere here an error occurs and bot remains unoperational after
	// queue has finished.
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

// Prev skips to previous song unconditionally. Means even if user has set
// repeat state to 'single' it will still skip to previous song. Kills current
// playback by sending to Stop channel.
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

// Current returns title of song with current SongIndex.
func (s *Stream) Current() string {
	if len(s.Queue) == 0 {
		return ""
	}
	return s.Queue[s.SongIndex].Title
}

// Add appends new song with given Source and Title to queue.
// REVIEW: how to add easier support for more fields
func (s *Stream) Add(Source, Title string) {
	s.Lock()
	defer s.Unlock()
	s.Queue = append(s.Queue, Song{Title, Source, len(s.Queue)})
}

// Skipto skips to song with given index.
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

// Disconnect resets stream and calls Disconnect method of current voice channel.
func (s *Stream) Disconnect() error {
	s.Reset()
	return s.V.Disconnect()
}

// Pause sets Playing state to false effectively pausing current ffmpeg process
// since the latter observes this flag.
func (s *Stream) Pause() {
	s.Playing = false
}

// Unpause sets playing flag to true effectively unpausing current ffmpeg
// process since the latter observes this flag.
func (s *Stream) Unpause() {
	s.Playing = true
}
