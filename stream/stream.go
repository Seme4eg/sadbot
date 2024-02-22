package stream

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/seme4eg/sadbot/track"
)

var ErrEmptyQueue = errors.New("empty queue")

type RepeatState string

const (
	RepeatOff    RepeatState = "off"
	RepeatSingle RepeatState = "single"
	RepeatAll    RepeatState = "all"
)

// Keys are guild IDs, 1 stream per guild.
type Streams map[string]*Stream

// Stream holds 1 stream per 1 guild, holds info about voice connection state,
// and other 'player' stuff like queue etc
type Stream struct {
	V       *discordgo.VoiceConnection
	Queue   []*track.Track
	current *track.Track
	repeat  RepeatState
}

// New returns new stream struct
func New(vc *discordgo.VoiceConnection) *Stream {
	return &Stream{
		V:      vc,
		repeat: RepeatOff,
	}
}

// Play starts playback of the CurrentTrack and runs next when playback finished
func (s *Stream) Play() error {
	if len(s.Queue) == 0 {
		return ErrEmptyQueue
	}

	if s.current == nil {
		s.current = s.Queue[0]
	}

	if s.current.IsPlaying() {
		return nil
	}

	c := make(chan error, 1)
	defer close(c)
	go func() {
		c <- s.current.Stream(s.V)
	}()

	if err, ok := <-c; ok {
		fmt.Println("error:", err)
		if err != nil {
			if err == track.ErrManualStop {
				fmt.Println("error manual")
				return nil
			}
			// REVIEW: maybe just return err and prettify it outside ?
			return fmt.Errorf("error playing audio file: %w", err)
		}
		// in case play function finished playing on its own (wasn't affected by
		// user commands) - skip to next track
		if s.repeat != RepeatSingle {
			return s.Next()
		} else {
			s.Play()
		}
	}

	return nil
}

// Reset resets fields of current Stream except session and stop channel.
// Stops possibly remaining ffmpeg process.
func (s *Stream) Reset() {
	s.Queue = s.Queue[:0]
	s.current.Stop()
	s.current = nil
	s.repeat = RepeatOff
}

// Clear empties queue and resets current track index to 0.
func (s *Stream) Clear() {
	s.Queue = s.Queue[:0]
}

// Shuffle shuffles queue. Currently playing track will be 1st always in shuffled
// queue.
func (s *Stream) Shuffle() {
	currentTrack := s.current
	// create temporary 'queue' value that doesn't contain currently playing
	// track since after each shuffle we want it to be still first in queue
	temp := append(s.Queue[:s.current.Index], s.Queue[s.current.Index+1:]...)
	rand.Shuffle(len(temp), func(i, j int) {
		temp[i], temp[j] = temp[j], temp[i]
	})
	s.Queue = append([]*track.Track{currentTrack}, temp...)
}

// Unshuffle sorts tracks in queue by their initial index field.
func (s *Stream) UnShuffle() {
	sort.Slice(s.Queue, func(i, j int) bool {
		return s.Queue[i].Index < s.Queue[j].Index
	})
}

// SetRepeat sets current guild's stream repeat state to either
// single / all or off.
func (s *Stream) SetRepeat(state string) error {
	switch val := RepeatState(state); val {
	case RepeatSingle, RepeatAll, RepeatOff:
		s.repeat = val
	default:
		return errors.New("invalid repeat state passed")
	}
	return nil
}

// Next skips to next track unconditionally. Means even if user has set repeat
// state to 'single' it will still skip to next track. Stops current playback.
func (s *Stream) Next() error {
	// FIXME: somewhere here an error occurs and bot remains unoperational after
	// queue has finished.
	if len(s.Queue) == 0 {
		return ErrEmptyQueue
	}
	index := s.current.Index + 1

	if index >= len(s.Queue) {
		switch s.repeat {
		case RepeatAll:
			index = 0
		case RepeatOff:
			// index = s.CurrentTrack.Index - 1
			return errors.New("no next track")
		}
	}

	s.current.Stop()
	s.current = s.Queue[index]
	s.Play()
	return nil
}

// Prev skips to previous track unconditionally. Means even if user has set
// repeat state to 'single' it will still skip to previous track. Kills current
// playback by sending to Stop channel.
func (s *Stream) Prev() error {
	index := s.current.Index - 1

	if index < 0 {
		switch s.repeat {
		case RepeatAll:
			index = len(s.Queue) - 1
		case RepeatOff:
			// index = s.CurrentTrack.Index + 1
			return errors.New("nothing was played before")
		}
	}

	s.current.Stop()
	s.current = s.Queue[index]
	s.Play()
	return nil
}

// Current returns title of current track
func (s *Stream) Current() string {
	if len(s.Queue) == 0 {
		return ""
	}
	return fmt.Sprint(s.current)
}

// Add adds new track with given source and title to queue.
func (s *Stream) Add(source, title string) {
	s.Queue = append(s.Queue, track.New(len(s.Queue), source, title))
}

// Skipto skips to track with given index.
func (s *Stream) Skipto(index int) error {
	if index <= 0 || index > len(s.Queue) {
		return errors.New("no track with such index")
	}

	s.current.Stop()
	s.current = s.Queue[index]

	return nil
}

// Disconnect resets stream and calls Disconnect method of current voice channel.
func (s *Stream) Disconnect() error {
	s.Reset()
	return s.V.Disconnect()
}

// Pause pauses currently playing track.
func (s *Stream) Pause() {
	s.current.Pause()
}

// Resume resumes currently playing track.
func (s *Stream) Resume() {
	s.current.Resume()
}
