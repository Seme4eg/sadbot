package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"sync"
)

// YtResult struct holds yt-dlp command output
type YtResult struct {
	Title    string  `json:"title"`
	Uploader string  `json:"uploader"`
	Duration float32 `json:"duration"`
	Url      string  `json:"webpage_url"`
}

// ProcessQuery determines the type of given query - urls or search.
// Gets information from yt-dlp on each url / search and returns as a slice
// of parsed results.
func ProcessQuery(query string) ([]YtResult, error) {
	t, err := QueryType(query)
	if err != nil {
		return []YtResult{}, err
	}

	switch t {
	case "search":
		tracks, err := Get(query)
		if err != nil {
			return []YtResult{}, err
		}
		return tracks, nil
	case "urls":
		// fetch info on each url in query (can be one url as well)
		res := []YtResult{}
		for _, url := range strings.Fields(query) {
			tracks, err := Get(url)
			if err != nil {
				fmt.Println("Failed to fetch", url)
				continue
			}
			res = append(res, tracks...)
		}
		return res, nil
	}
	return []YtResult{}, errors.New("something went wrong when determining query type")
}

// QueryType checks if all whitespace-separated strings are either urls or just
// strings returns false in case there's a mix of both.
func QueryType(query string) (string, error) {
	allUrls := true
	allWords := true
	for _, u := range strings.Fields(query) {
		if IsUrl(u) {
			allWords = false
		} else {
			allUrls = false
		}
	}
	if allUrls {
		return "urls", nil
	}
	if allWords {
		return "search", nil
	}
	return "", errors.New("either all arguments must be urls or none")
}

// IsUrl determines if given string is a url.
func IsUrl(link string) bool {
	parsedURL, err := url.ParseRequestURI(link)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	} else {
		return true
	}
}

// Get calls yt-dlp and gets information on provided query, which can also be a
// playlist, which is why it returns a slice
func Get(query string) ([]YtResult, error) {
	// default search to "ytsearch" to keep things simple
	// as an option - in the future allow user to specify search method
	// like "ytsearch : search query" maybe
	ytDlp := exec.Command("yt-dlp", "--flat-playlist", "--default-search",
		"ytsearch", "--print", "%(.{title,uploader,duration,webpage_url})j", query)
	var out bytes.Buffer
	ytDlp.Stdout = &out
	err := ytDlp.Run()
	if err != nil {
		return []YtResult{}, err
	}

	tracks := []YtResult{}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		track := YtResult{}
		err = json.Unmarshal(scanner.Bytes(), &track)
		if err != nil {
			return []YtResult{}, err
		}
		tracks = append(tracks, track)
	}

	var wg sync.WaitGroup

	// FIXME: if no title in output then prob given link was soundcloud playlist
	// then we need to fetch info track by track. Maybe there is a way to fetch
	// needed playlist info with only 1 command, but for now i haven't found how.
	for i, t := range tracks {
		track := t
		index := i
		if track.Title == "" {
			wg.Add(1)

			go func() {
				ytDlp := exec.Command("yt-dlp", "--flat-playlist", "--default-search",
					"ytsearch", "--print", "%(.{title,uploader,duration,webpage_url})j", track.Url)
				var out bytes.Buffer
				ytDlp.Stdout = &out
				err := ytDlp.Run()
				if err != nil {
					fmt.Printf("Error processing %s : %s", track.Url, err)
				}
				err = json.Unmarshal(out.Bytes(), &track)
				if err != nil {
					fmt.Printf("Error unmarshaling response for %s : %s", track.Url, err)
				}
				if track.Uploader != "" {
					track.Title = track.Uploader + " – " + track.Title
				}
				tracks[index] = track
				defer wg.Done()
			}()
		}
	}

	wg.Wait()

	return tracks, nil
}
