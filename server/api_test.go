package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/mock"
)

// goneWriter behaves like the plugin RPC writer of a request the browser has
// cancelled: its header map is nil, so any http.Error on it panics.
type goneWriter struct{ wrote bool }

func (g *goneWriter) Header() http.Header { return nil }
func (g *goneWriter) Write(b []byte) (int, error) {
	g.wrote = true
	return len(b), nil
}
func (g *goneWriter) WriteHeader(int) { g.wrote = true }

func testPlugin(t *testing.T) (*Plugin, *plugintest.API) {
	t.Helper()
	api := &plugintest.API{}
	api.On("GetUser", mock.Anything).Return(&model.User{Locale: "es"}, nil).Maybe()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	p := &Plugin{tracks: newTrackCache()}
	p.SetAPI(api)
	p.yt, _ = fixtureServer(t)
	p.router = p.newRouter()
	return p, api
}

// Typing in the search window cancels the previous search; that used to crash
// the plugin process (18/11/2026).
func TestSearchCancelledDoesNotAnswer(t *testing.T) {
	p, api := testPlugin(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/search?q=rosalia", nil)
	r.Header.Set("Mattermost-User-Id", "u")
	w := &goneWriter{}
	p.router.ServeHTTP(w, r)
	if w.wrote {
		t.Error("a cancelled search must not be answered")
	}
	api.AssertNotCalled(t, "LogWarn", mock.Anything, mock.Anything, mock.Anything)
	api.AssertNotCalled(t, "LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSearchAnswers(t *testing.T) {
	p, _ := testPlugin(t)
	r := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=rosalia", nil)
	r.Header.Set("Mattermost-User-Id", "u")
	w := httptest.NewRecorder()
	p.router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d: %s", w.Code, w.Body)
	}
	if _, ok := p.tracks.get("dN3y8kVNplA"); !ok {
		t.Error("search results must be cached for sharing")
	}
}

func TestRecoverPanics(t *testing.T) {
	p, api := testPlugin(t)
	h := p.recoverPanics(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		var m map[string]int
		m["x"] = 1
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/x", nil))
	api.AssertCalled(t, "LogError", "Panic serving a request", "path", "/x", "panic", mock.Anything, "stack", mock.Anything)
}
