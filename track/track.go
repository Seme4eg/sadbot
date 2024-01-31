package track

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// REVIEW: how to add easier support for more fields
type Track struct {
	title string
	// source can be either Url or a path (if using playfolder command)
	source  string
	Index   int // initial index, used to 'unshuffle' suffled queue
	playing bool
	// TODO: Duration string
	// ctx is being passed to play func process
	ctx context.Context
	// used on stop / next / prev functions from stream and cancelles currently
	// processing playfunc
	// REVIEW: whether cancel func should be in 'track' itself or in in 'stream'
	cancel context.CancelFunc
}

func New(i int, source, title string) *Track {
	ctx, cancel := context.WithCancel(context.Background())
	return &Track{
		title:  title,
		source: source,
		Index:  i,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *Track) Play(v *discordgo.VoiceConnection) (err error) {
	t.playing = true
	c := make(chan error, 1)

	// start playback function in background
	go func() {
		c <- playAudioFile(t.ctx, t, v)
	}()

	// Send "speaking" packet over the voice websocket
	if err = v.Speaking(true); err != nil {
		return err
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			log.Println("Couldn't stop speaking:", err)
		}
	}()

	select {
	case <-t.ctx.Done():
		t.playing = false
		return t.ctx.Err()
	case res := <-c:
		t.playing = false
		return res
	}
}

func (t *Track) Done() <-chan struct{} {
	// don't expose anything of 'ctx' but 'Done()' method
	return t.ctx.Done()
}

func (t *Track) Cancel() {
	t.playing = false
	t.cancel()
}

// Pause sets Playing state to false, which pauses current ffmpeg process.
func (t *Track) Pause() {
	t.playing = false
}

// Resume sets Playing state to true, which continues current ffmpeg process.
func (t *Track) Resume() {
	t.playing = false
}

func (t *Track) IsPlaying() bool {
	return t.playing
}

func (t *Track) String() string {
	format := "%v"
	if t.playing {
		format = "**%v**"
	}
	return fmt.Sprintf(format, t.title)
}
