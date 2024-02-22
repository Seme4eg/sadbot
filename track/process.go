package track

import (
	"io"
	"os/exec"
	"strconv"

	"github.com/seme4eg/sadbot/utils"
)

type Process interface {
	Start() error
	Kill() error
	StdoutPipe() (io.ReadCloser, error)
}

func NewProcess(source string) (Process, error) {
	if utils.IsUrl(source) {
		return NewUrlProcess(source)
	} else {
		return NewFileProcess(source), nil
	}
}

type FileProcess struct {
	ffmpeg *exec.Cmd
}

func NewFileProcess(source string) *FileProcess {
	ffmpeg := exec.Command("ffmpeg", "-i", source, "-f", "s16le", "-ar",
		strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	return &FileProcess{ffmpeg}
}

func (p *FileProcess) Start() error {
	return p.ffmpeg.Start()
}

func (p *FileProcess) Kill() error {
	return p.ffmpeg.Process.Kill()
}

func (p *FileProcess) StdoutPipe() (io.ReadCloser, error) {
	return p.ffmpeg.StdoutPipe()
}

type UrlProcess struct {
	ffmpeg *exec.Cmd
	ytdlp  *exec.Cmd
}

func NewUrlProcess(source string) (*UrlProcess, error) {
	ytdlp := exec.Command("yt-dlp", "--no-part", "--downloader", "ffmpeg",
		"--buffer-size", "16K", "--limit-rate", "50K", "-o", "-", "-f", "bestaudio", source)

	ytdlpOut, err := ytdlp.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// FIXME: still sometimes skips to next song before current finished playing
	// Prevent yt-dlp command to finish before ffmpeg is done reading its output
	// go func() {
	// 	if err := ytdlp.Wait(); err != nil {
	// 		fmt.Println("error waiting for ytdlp to finish:", err)
	// 	}
	// }()

	ffmpeg := exec.Command("ffmpeg", "-i", "-", "-f", "s16le", "-ar",
		strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	ffmpeg.Stdin = ytdlpOut

	return &UrlProcess{ffmpeg, ytdlp}, nil
}

func (p *UrlProcess) Start() error {
	if err := p.ytdlp.Start(); err != nil {
		return err
	}
	return p.ffmpeg.Start()
}

func (p *UrlProcess) Kill() error {
	if err := p.ffmpeg.Process.Kill(); err != nil {
		return err
	}
	return p.ytdlp.Process.Kill()
}

func (p *UrlProcess) StdoutPipe() (io.ReadCloser, error) {
	return p.ffmpeg.StdoutPipe()
}
