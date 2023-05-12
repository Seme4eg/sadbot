// credits to https://github.com/bwmarrin/dgvoice/blob/master/dgvoice.go
// i only 'obfuscated' this file a bit

package utils

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

// values below seem to provide the best overall performance, others might b
// unstable
const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	// mu          sync.Mutex
)

type AudioError string

func (a AudioError) Error() string {
	return "Audio module error: " + string(a)
}

// SendPCM receives on the provied channel
// encodes received PCM data into Opus and sends that to Discordgo
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) error {
	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		return AudioError("NewEncoder Error" + err.Error())
	}

	for {
		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			return AudioError("PCM Channel closed")
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return AudioError("Encoding Error" + err.Error())
		}

		if !v.Ready || v.OpusSend == nil {
			// Sending errors here might not be suited
			return nil
		}

		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}

// PlayAudioFile will play the given filename to the already connected
// Discord voice server/channel. voice websocket and udp socket
// must already be setup before this will work.
func PlayAudioFile(v *discordgo.VoiceConnection, filename string, stop <-chan bool, playing *bool) error {
	// create ffmpeg command
	ffmpeg := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar",
		strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return AudioError("StdoutPipe Error" + err.Error())
	}

	// read in chunks of 16KB (16 / 1024 bytes)
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// Starts the ffmpeg command
	if err = ffmpeg.Start(); err != nil {
		return AudioError("RunStart Error" + err.Error())
	}

	// prevent memory leak from residual ffmpeg streams
	defer ffmpeg.Process.Kill()

	// when stop is sent, kill ffmpeg
	go func() {
		<-stop
		err = ffmpeg.Process.Kill()
		if err != nil {
			fmt.Println("Error killing ffmpeg process:", err)
		}
	}()

	// Send "speaking" packet over the voice websocket
	if err = v.Speaking(true); err != nil {
		return AudioError("Couldn't set speaking" + err.Error())
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			fmt.Println("Couldn't stop speaking:", err)
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		err := SendPCM(v, send)
		if err != nil {
			fmt.Println("SendPCM error:", err)
		}
		close <- true
	}()

	for {
		// means player was paused by the user, check every second on status change
		if !*playing {
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
			return AudioError("error reading from ffmpeg stdout" + err.Error())
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return nil
		}
	}
}
