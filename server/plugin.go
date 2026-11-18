package main

import (
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"
)

// The card's song is in the post props; its text is a plain link for clients
// that do not run the webapp plugin.
const postType = "custom_music_card"

type Plugin struct {
	plugin.MattermostPlugin
	router   http.Handler
	settings settings
	yt       *ytMusic
	tracks   *trackCache
}

func (p *Plugin) OnActivate() error {
	p.yt = newYTMusic(p.settings.region)
	p.tracks = newTrackCache()
	p.router = p.newRouter()
	return nil
}

func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func main() {
	plugin.ClientMain(&Plugin{})
}
