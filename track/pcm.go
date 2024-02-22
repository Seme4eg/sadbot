package track

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

var (
	opusEncoder  *gopus.Encoder
	ErrPcmClosed = errors.New("err: PCM Channel closed")
)

// sendPCM receives on the provied channel
// encodes received PCM data into Opus and sends that to Discordgo
func sendPCM(stop <-chan struct{}, v *discordgo.VoiceConnection, pcm <-chan []int16) error {
	// Send "speaking" packet over the voice websocket
	if err := v.Speaking(true); err != nil {
		return err
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			log.Println("Couldn't stop speaking:", err)
		}
	}()

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		return err
	}

	for {
		select {
		case <-stop:
			fmt.Println("pcm stopped")

			return nil
		case recv, ok := <-pcm:
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
}
