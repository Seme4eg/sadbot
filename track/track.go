package track

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Track struct {
	title string
	// source can be either Url or a path (if using playfolder command)
	source string
	// REVIEW: remove, track doesn't know about stream ?
	Index   int // initial index, used to 'unshuffle' suffled queue
	playing bool
	stop    chan struct{}
	// TODO: Duration string
}

var ErrManualStop = errors.New("manual process stop")

func New(i int, source, title string) *Track {
	return &Track{
		title:  title,
		source: source,
		Index:  i,
	}
}

func (t *Track) Stream(con *discordgo.VoiceConnection) (err error) {
	t.playing = true

	c := make(chan error, 1)
	t.stop = make(chan struct{})
	// don't make 2 separate channels for read and write, will result in memory
	// pollution on large scale. Keep those funciton in sync
	audioChan := make(chan []int16, 50)

	defer func() {
		// close(c)
		// close(t.stop)
		// close(audioChan)
		t.playing = false
	}()

	go func() {
		c <- t.readAudioFile(audioChan)
	}()

	go func() {
		c <- sendPCM(t.stop, con, audioChan)
	}()

	select {
	case <-t.stop:
		// FIXME:
		// this case can happen only on user interaction, means we do not need
		// to respect current Repeat state
		// prevent ffmpeg processes from overlapping (especially on prev command)
		// time.Sleep(time.Millisecond * 350)
		return ErrManualStop
	case res := <-c:
		if errors.Is(err, ErrPcmClosed) {
			return nil
		}
		return res
	}
}

// readAudioFile plays the given filename to the already connected Discord voice
// server/channel. Voice websocket and udp socket must already be setup before
// this will work.
// credits to https://github.com/bwmarrin/dgvoice/blob/master/dgvoice.go
// for hints
func (t *Track) readAudioFile(audiochan chan<- []int16) error {
	readProcess, err := NewProcess(t.source)
	if err != nil {
		return err
	}

	ffmpegout, err := readProcess.StdoutPipe()
	if err != nil {
		return err
	}

	// read in chunks of 16KB (16 / 1024 bytes)
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	if err = readProcess.Start(); err != nil {
		return err
	}

	// prevent memory leak from residual ffmpeg streams
	// NOTE: not handling errors since the only errors that can appear here is
	// that the process no longer exists
	defer readProcess.Kill()

	for {
		select {
		case <-t.stop:
			fmt.Println("ffmpeg stopped")
			err := readProcess.Kill()
			if err != nil {
				fmt.Println("Error killing ffmpeg process:", err)
			}
			return nil
		default:
			// means player was paused by the user, check every second on status change
			if !t.playing {
				time.Sleep(1 * time.Second)
				continue
			}
			// read data from ffmpeg stdout
			audiobuf := make([]int16, frameSize*channels)
			err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
			if err != nil {
				return err
			}

			audiochan <- audiobuf
		}
	}
}

// Pause sets Playing state to false, which pauses current ffmpeg process.
func (t *Track) Pause() {
	t.playing = false
}

// Resume sets Playing state to true, which continues current ffmpeg process.
func (t *Track) Resume() {
	t.playing = true
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

func (t *Track) Stop() {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	if t.stop != nil {
		fmt.Println("sending to stop chan")
		t.stop <- struct{}{}
	}
}
