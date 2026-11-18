package main

import (
	"net/url"
	"regexp"
	"strings"
)

var videoIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

func validVideoID(id string) bool {
	return videoIDRe.MatchString(id)
}

// videoIDFromLink returns the video of a YouTube or YouTube Music link, so a
// pasted link becomes a card without searching.
func videoIDFromLink(input string) (string, bool) {
	u, err := url.Parse(strings.TrimSpace(input))
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return "", false
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	var id string
	switch host {
	case "youtu.be":
		id = strings.Trim(u.Path, "/")
	case "youtube.com", "m.youtube.com", "music.youtube.com":
		if u.Path == "/watch" {
			id = u.Query().Get("v")
		} else if rest, ok := strings.CutPrefix(u.Path, "/shorts/"); ok {
			id = rest
		}
	}
	return id, validVideoID(id)
}
