package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// YouTube Music's own web API (InnerTube), the one music.youtube.com calls. It
// needs no key and has no quota, but it is not documented: if Google changes
// it, search breaks while pasted links keep working through lookup.
const (
	innertubeURL = "https://music.youtube.com/youtubei/v1/"
	clientName   = "WEB_REMIX"
	// The "Songs" filter of the search page.
	songsFilter = "EgWKAQIIAWoMEA4QChADEAQQCRAF"
	maxResults  = 10
	// Covers are served at any size by changing the =wN-hN suffix.
	coverSize = 544
)

var (
	durationRe = regexp.MustCompile(`^\d+(:\d{2}){1,2}$`)
	sizeRe     = regexp.MustCompile(`=w\d+-h\d+`)
)

type track struct {
	VideoID  string `json:"video_id"`
	Title    string `json:"title"`
	Artists  string `json:"artists"`
	Album    string `json:"album,omitempty"`
	Duration string `json:"duration,omitempty"`
	Cover    string `json:"cover,omitempty"`
}

type ytMusic struct {
	http    *http.Client
	baseURL string
	region  func() string
}

func newYTMusic(region func() string) *ytMusic {
	return &ytMusic{http: &http.Client{Timeout: 10 * time.Second}, baseURL: innertubeURL, region: region}
}

// clientVersion follows the web client's date-based scheme; the API refuses
// versions that are too old, so it is built from today's date like ytmusicapi does.
func clientVersion(now time.Time) string {
	return "1." + now.UTC().Format("20060102") + ".01.00"
}

func (y *ytMusic) call(ctx context.Context, endpoint, lang string, body map[string]any, out any) error {
	body["context"] = map[string]any{"client": map[string]string{
		"clientName":    clientName,
		"clientVersion": clientVersion(time.Now()),
		"hl":            lang,
		"gl":            y.region(),
	}}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, y.baseURL+endpoint+"?prettyPrint=false", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://music.youtube.com")
	resp, err := y.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("youtube music %s: HTTP %d", endpoint, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (y *ytMusic) search(ctx context.Context, query, lang string) ([]track, error) {
	var resp searchResponse
	if err := y.call(ctx, "search", lang, map[string]any{"query": query, "params": songsFilter}, &resp); err != nil {
		return nil, err
	}
	return resp.tracks(), nil
}

func (y *ytMusic) lookup(ctx context.Context, videoID, lang string) (track, error) {
	var resp playerResponse
	if err := y.call(ctx, "player", lang, map[string]any{"videoId": videoID}, &resp); err != nil {
		return track{}, err
	}
	return resp.track()
}

type text struct {
	Runs []run `json:"runs"`
}

type run struct {
	Text               string `json:"text"`
	NavigationEndpoint struct {
		BrowseEndpoint struct {
			Configs struct {
				Music struct {
					PageType string `json:"pageType"`
				} `json:"browseEndpointContextMusicConfig"`
			} `json:"browseEndpointContextSupportedConfigs"`
		} `json:"browseEndpoint"`
	} `json:"navigationEndpoint"`
}

type thumbnails struct {
	Thumbnails []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
}

type listItem struct {
	Thumbnail struct {
		Renderer struct {
			Thumbnail thumbnails `json:"thumbnail"`
		} `json:"musicThumbnailRenderer"`
	} `json:"thumbnail"`
	FlexColumns []struct {
		Column struct {
			Text text `json:"text"`
		} `json:"musicResponsiveListItemFlexColumnRenderer"`
	} `json:"flexColumns"`
	PlaylistItemData struct {
		VideoID string `json:"videoId"`
	} `json:"playlistItemData"`
}

type searchResponse struct {
	Contents struct {
		Tabbed struct {
			Tabs []struct {
				Tab struct {
					Content struct {
						SectionList struct {
							Contents []struct {
								Shelf *struct {
									Contents []struct {
										Item *listItem `json:"musicResponsiveListItemRenderer"`
									} `json:"contents"`
								} `json:"musicShelfRenderer"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"tabbedSearchResultsRenderer"`
	} `json:"contents"`
}

func (r searchResponse) tracks() []track {
	tracks := []track{}
	for _, tab := range r.Contents.Tabbed.Tabs {
		for _, section := range tab.Tab.Content.SectionList.Contents {
			if section.Shelf == nil {
				continue
			}
			for _, entry := range section.Shelf.Contents {
				if t, ok := entry.Item.track(); ok {
					tracks = append(tracks, t)
				}
				if len(tracks) == maxResults {
					return tracks
				}
			}
		}
	}
	return tracks
}

// track reads a search row: the title in the first column and
// "Artist y Artist • Album • 2:37" in the second.
func (item *listItem) track() (track, bool) {
	if item == nil || !validVideoID(item.PlaylistItemData.VideoID) || len(item.FlexColumns) < 2 {
		return track{}, false
	}
	t := track{VideoID: item.PlaylistItemData.VideoID, Title: joinRuns(item.FlexColumns[0].Column.Text.Runs)}
	var groups [][]run
	current := []run{}
	for _, r := range item.FlexColumns[1].Column.Text.Runs {
		if strings.TrimSpace(r.Text) == "•" {
			groups = append(groups, current)
			current = []run{}
			continue
		}
		current = append(current, r)
	}
	groups = append(groups, current)
	t.Artists = joinRuns(groups[0])
	for _, group := range groups[1:] {
		value := joinRuns(group)
		switch {
		case durationRe.MatchString(value):
			t.Duration = value
		case len(group) > 0 && group[0].NavigationEndpoint.BrowseEndpoint.Configs.Music.PageType == "MUSIC_PAGE_TYPE_ALBUM":
			t.Album = value
		}
	}
	t.Cover = coverURL(item.Thumbnail.Renderer.Thumbnail)
	return t, t.Title != ""
}

type playerResponse struct {
	VideoDetails struct {
		VideoID       string     `json:"videoId"`
		Title         string     `json:"title"`
		Author        string     `json:"author"`
		LengthSeconds string     `json:"lengthSeconds"`
		Thumbnail     thumbnails `json:"thumbnail"`
	} `json:"videoDetails"`
}

// track reads the details of a single video. The player itself says
// "UNPLAYABLE" without a browser token, but the details are still there.
func (r playerResponse) track() (track, error) {
	d := r.VideoDetails
	if !validVideoID(d.VideoID) || d.Title == "" {
		return track{}, fmt.Errorf("video not found")
	}
	t := track{
		VideoID: d.VideoID,
		Title:   d.Title,
		// Auto-generated artist channels are called "Artist - Topic".
		Artists: strings.TrimSuffix(d.Author, " - Topic"),
		Cover:   coverURL(d.Thumbnail),
	}
	if seconds, err := strconv.Atoi(d.LengthSeconds); err == nil && seconds > 0 {
		t.Duration = formatDuration(seconds)
	}
	return t, nil
}

func joinRuns(runs []run) string {
	var b strings.Builder
	for _, r := range runs {
		b.WriteString(r.Text)
	}
	return strings.TrimSpace(b.String())
}

func coverURL(t thumbnails) string {
	if len(t.Thumbnails) == 0 {
		return ""
	}
	url := t.Thumbnails[len(t.Thumbnails)-1].URL
	return sizeRe.ReplaceAllString(url, fmt.Sprintf("=w%d-h%d", coverSize, coverSize))
}

func formatDuration(seconds int) string {
	h, m, s := seconds/3600, seconds/60%60, seconds%60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}
