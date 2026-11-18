package main

import "sync"

// Searches return the album, which the single-video lookup does not; keeping
// recent results lets the shared card show it without asking YouTube again.
const maxCachedTracks = 1000

type trackCache struct {
	mu     sync.Mutex
	tracks map[string]track
}

func newTrackCache() *trackCache {
	return &trackCache{tracks: map[string]track{}}
}

func (c *trackCache) add(tracks []track) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Starting over when full is enough: a miss only costs one lookup.
	if len(c.tracks)+len(tracks) > maxCachedTracks {
		c.tracks = map[string]track{}
	}
	for _, t := range tracks {
		c.tracks[t.VideoID] = t
	}
}

func (c *trackCache) get(videoID string) (track, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, ok := c.tracks[videoID]
	return t, ok
}
