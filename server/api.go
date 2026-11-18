package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/mattermost/mattermost/server/public/model"
)

const maxQueryLength = 200

func (p *Plugin) newRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search", p.handleSearch)
	mux.HandleFunc("POST /api/v1/share", p.handleShare)
	return requireUser(mux)
}

func requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Mattermost-User-Id") == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(v)
}

// handleSearch returns songs for a text, or the one song of a pasted link.
func (p *Plugin) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" || utf8.RuneCountInString(query) > maxQueryLength {
		http.Error(w, "invalid query", http.StatusBadRequest)
		return
	}
	lang := p.userLanguage(r.Header.Get("Mattermost-User-Id"))
	var tracks []track
	if id, ok := videoIDFromLink(query); ok {
		t, err := p.yt.lookup(r.Context(), id, lang)
		if err != nil {
			p.API.LogWarn("YouTube Music lookup failed", "video", id, "err", err.Error())
			http.Error(w, "song not found", http.StatusBadGateway)
			return
		}
		tracks = []track{t}
	} else {
		var err error
		if tracks, err = p.yt.search(r.Context(), query, lang); err != nil {
			p.API.LogWarn("YouTube Music search failed", "err", err.Error())
			http.Error(w, "search failed", http.StatusBadGateway)
			return
		}
	}
	p.tracks.add(tracks)
	writeJSON(w, map[string][]track{"tracks": tracks})
}

type shareRequest struct {
	ChannelID string `json:"channel_id"`
	RootID    string `json:"root_id"`
	VideoID   string `json:"video_id"`
	Text      string `json:"text"`
}

// handleShare posts the card. The song details come from the server (its
// recent searches, or YouTube Music), never from the request, so a card always
// shows what the video really is.
func (p *Plugin) handleShare(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	var req shareRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024)).Decode(&req); err != nil || !validVideoID(req.VideoID) {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(req.Text) > model.PostMessageMaxRunesV2/2 {
		http.Error(w, "text too long", http.StatusRequestEntityTooLarge)
		return
	}
	if err := p.checkCanPost(userID, req.ChannelID, req.RootID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	t, ok := p.tracks.get(req.VideoID)
	if !ok {
		var err error
		if t, err = p.yt.lookup(r.Context(), req.VideoID, p.userLanguage(userID)); err != nil {
			http.Error(w, "song not found", http.StatusBadGateway)
			return
		}
	}
	post, appErr := p.API.CreatePost(&model.Post{
		UserId:    userID,
		ChannelId: req.ChannelID,
		RootId:    req.RootID,
		Type:      postType,
		Message:   postMessage(req.Text, t),
		Props:     trackProps(req.Text, t),
	})
	if appErr != nil {
		p.API.LogError("Failed to post a music card", "err", appErr.Error())
		http.Error(w, "failed to post", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"post_id": post.Id})
}

func (p *Plugin) checkCanPost(userID, channelID, rootID string) error {
	channel, appErr := p.API.GetChannel(channelID)
	if appErr != nil || channel.DeleteAt != 0 {
		return errors.New("channel not found")
	}
	if !p.API.HasPermissionToChannel(userID, channelID, model.PermissionCreatePost) {
		return errors.New("you cannot post in this channel")
	}
	if rootID == "" {
		return nil
	}
	root, appErr := p.API.GetPost(rootID)
	if appErr != nil || root.ChannelId != channelID || root.RootId != "" {
		return errors.New("thread not found")
	}
	return nil
}

func (p *Plugin) userLanguage(userID string) string {
	if user, appErr := p.API.GetUser(userID); appErr == nil && user.Locale != "" {
		return strings.Split(user.Locale, "-")[0]
	}
	return "en"
}

func musicURL(videoID string) string {
	return "https://music.youtube.com/watch?v=" + videoID
}

// postMessage is what notifications, search and the mobile apps (which do not
// run webapp plugins) show: the user's text and a link to the song.
func postMessage(text string, t track) string {
	link := fmt.Sprintf("🎵 [%s — %s](%s)", escapeLinkText(t.Title), escapeLinkText(t.Artists), musicURL(t.VideoID))
	if text = strings.TrimSpace(text); text != "" {
		return text + "\n\n" + link
	}
	return link
}

// escapeLinkText keeps titles like "Song [Live]" from breaking the markdown link.
var escapeLinkText = strings.NewReplacer(`\`, `\\`, `[`, `\[`, `]`, `\]`).Replace

func trackProps(text string, t track) model.StringInterface {
	return model.StringInterface{
		"caption":  strings.TrimSpace(text),
		"video_id": t.VideoID,
		"title":    t.Title,
		"artists":  t.Artists,
		"album":    t.Album,
		"duration": t.Duration,
		"cover":    t.Cover,
	}
}
