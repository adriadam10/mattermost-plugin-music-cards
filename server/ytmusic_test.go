package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// The fixtures are trimmed real responses of music.youtube.com (November 2026).
func fixtureServer(t *testing.T) (*ytMusic, *map[string]any) {
	t.Helper()
	var last map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&last)
		name := map[string]string{"/search": "testdata/search.json", "/player": "testdata/player.json"}[r.URL.Path]
		body, err := os.ReadFile(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	y := newYTMusic(func() string { return "ES" })
	y.baseURL = srv.URL + "/"
	return y, &last
}

func TestSearch(t *testing.T) {
	y, last := fixtureServer(t)
	tracks, err := y.search(context.Background(), "rosalia despecha", "es")
	if err != nil {
		t.Fatal(err)
	}
	want := []track{
		{VideoID: "dN3y8kVNplA", Title: "DESPECHÁ", Artists: "ROSALÍA", Album: "MOTOMAMI +", Duration: "2:37"},
		{VideoID: "w7pjt9ZH3NM", Title: "La Perla", Artists: "ROSALÍA y Yahritza Y Su Esencia", Album: "LUX", Duration: "3:16"},
	}
	if len(tracks) != len(want) {
		t.Fatalf("got %d tracks: %+v", len(tracks), tracks)
	}
	for i, w := range want {
		got := tracks[i]
		if got.Cover == "" || !sizeRe.MatchString(got.Cover) {
			t.Errorf("track %d cover %q", i, got.Cover)
		}
		got.Cover = ""
		if got != w {
			t.Errorf("track %d: got %+v, want %+v", i, got, w)
		}
	}
	client := (*last)["context"].(map[string]any)["client"].(map[string]any)
	if client["gl"] != "ES" || client["hl"] != "es" || (*last)["params"] != songsFilter {
		t.Errorf("request %+v", *last)
	}
}

func TestLookup(t *testing.T) {
	y, _ := fixtureServer(t)
	got, err := y.lookup(context.Background(), "w7pjt9ZH3NM", "es")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "La Perla" || got.Artists != "ROSALÍA y Yahritza Y Su Esencia" || got.Duration != "3:16" {
		t.Errorf("got %+v", got)
	}
	if want := "=w544-h544"; !strings.Contains(got.Cover, want) {
		t.Errorf("cover %q lacks %s", got.Cover, want)
	}
}

func TestLookupMissing(t *testing.T) {
	if _, err := (playerResponse{}).track(); err == nil {
		t.Error("an empty player response must fail")
	}
}

func TestFormatDuration(t *testing.T) {
	for in, want := range map[int]string{5: "0:05", 157: "2:37", 3725: "1:02:05"} {
		if got := formatDuration(in); got != want {
			t.Errorf("%d: got %s, want %s", in, got, want)
		}
	}
}

func TestClientVersion(t *testing.T) {
	if got := clientVersion(time.Date(2026, 11, 18, 23, 0, 0, 0, time.UTC)); got != "1.20261118.01.00" {
		t.Errorf("got %s", got)
	}
}
