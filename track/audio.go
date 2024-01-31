// credits to https://github.com/bwmarrin/dgvoice/blob/master/dgvoice.go
// i only added urls support, piping from yt-dlp if passed string was url,
// pause support and obfuscated this file a bit. All changes can be seen in this
// file git history.

package track

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/seme4eg/sadbot/utils"
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
	speakers     map[uint32]*gopus.Decoder
	opusEncoder  *gopus.Encoder
	ErrPcmClosed = errors.New("err: PCM Channel closed")
)

// sendPCM receives on the provied channel
// encodes received PCM data into Opus and sends that to Discordgo
func sendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) error {
	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		return err
	}

	for {
		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			return ErrPcmClosed
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return err
		}

		if !v.Ready || v.OpusSend == nil {
			// Sending errors here might not be suited
			return nil
		}

		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}

// playAudioFile plays the given filename to the already connected Discord voice
// server/channel. Voice websocket and udp socket must already be setup before
// this will work.
// FIXME: split this huge function
func playAudioFile(ctx context.Context, track *Track, v *discordgo.VoiceConnection) (err error) {

	ffmpeg, err := createFfmpegProcess(track.source)
	if err != nil {
		return fmt.Errorf("error starting ffmpeg process: %s", err)
	}

	ffmpegout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}

	// read in chunks of 16KB (16 / 1024 bytes)
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	if err = ffmpeg.Start(); err != nil {
		return err
	}

	// prevent memory leak from residual ffmpeg streams
	// NOTE: not handling errors since the only errors that can appear here is
	// that the process no longer exists
	defer ffmpeg.Process.Kill()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		err := sendPCM(v, send)
		if err != nil {
			// ignore pcm closed error since it appears on every process end
			if !errors.Is(err, ErrPcmClosed) {
				fmt.Println("SendPCM error:", err)
			}
		}
		close <- true
	}()

	for {
		// means player was paused by the user, check every second on status change
		if !track.playing {
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

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return nil
		}
	}
}

func createFfmpegProcess(source string) (ffmpeg *exec.Cmd, err error) {
	var (
		ytdlp    *exec.Cmd
		ytdlpOut io.ReadCloser
	)

	isUrl := utils.IsUrl(source)

	if isUrl {
		ytdlp = exec.Command("yt-dlp", "--no-part", "--downloader", "ffmpeg",
			"--buffer-size", "16K", "--limit-rate", "50K", "-o", "-", "-f", "bestaudio", source)
		// since ytdlp is now source of ffmpeg command we need to change source
		// to "-" so ffmpeg reads from pipe
		source = "-"
		ytdlpOut, err = ytdlp.StdoutPipe()
		if err != nil {
			return nil, err
		}
		if err := ytdlp.Start(); err != nil {
			return nil, err
		}

		// FIXME: still sometimes skips to next song before current finished playing
		// Prevent yt-dlp command to finish before ffmpeg is done reading its output
		// go func() {
		// 	if err := ytdlp.Wait(); err != nil {
		// 		fmt.Println("error waiting for ytdlp to finish:", err)
		// 	}
		// }()

		// NOTE: decided not to handle errors on ytdlp kill for now
		defer ytdlp.Process.Kill()
	}

	ffmpeg = exec.Command("ffmpeg", "-i", source, "-f", "s16le", "-ar",
		strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	if isUrl {
		ffmpeg.Stdin = ytdlpOut
	}

	return ffmpeg, nil
}
