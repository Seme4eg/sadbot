package utils

import ()

type Song struct {
	Media    string
	Title    string
	Duration *string
	Id       string
}

// func (song Song) Ffmpeg() *exec.Cmd {
// 	return exec.Command("ffmpeg", "-i", song.Media, "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac",
// 		strconv.Itoa(CHANNELS), "pipe:1")
// }
