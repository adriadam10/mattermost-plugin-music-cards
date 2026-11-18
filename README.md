# Mattermost Music Cards plugin

Share songs from YouTube Music as cards with a built-in player.

- A music-note button in the message box toolbar opens a search window. Type a song,
  artist or album and pick one of the results (mouse, or ↑ ↓ and Enter). Pasting a
  YouTube or YouTube Music link works too.
- Whatever is written in the message box goes along as the card's text (markdown,
  mentions and emoji included) and the box is cleared.
- The card shows the cover, title, artists, album and length. Pressing the cover
  opens the YouTube player inside the card; starting another card stops the first one.
  The title and the ↗ button open the song in YouTube Music.
- Works in direct messages, group messages, channels and threads.

## How it works

1. The search window calls `GET /api/v1/search?q=…`. The server asks YouTube Music's
   own web API (the one music.youtube.com uses, with the *Songs* filter) and returns
   up to 10 songs. A pasted link is looked up instead of searched.
2. Choosing a song calls `POST /api/v1/share` with the channel, the thread, the video
   id and the text. The server takes the song details from its recent searches (or
   looks the video up again), never from the request, and creates a post of type
   `custom_music_card` with them in its props.
3. The post's message is the text plus a `🎵 [Title — Artists](link)` line, which is what
   notifications, search and the mobile apps show.

The player is `youtube-nocookie.com` and only loads when someone presses play, so
reading a channel does not contact YouTube except for the cover images.

## Settings

System Console → Plugins → Music Cards:

- **Region**: two-letter country code used for search results and availability
  (default `US`).

## Limitations

- YouTube Music's web API needs no key and has no quota, but it is not documented.
  If Google changes it, search stops working until the plugin is updated; pasted
  links keep working through the lookup.
- Some labels do not allow their videos to play outside YouTube: the player then
  shows YouTube's own message and the ↗ button still opens the song.
- The native mobile apps do not run webapp plugins: they show the text and the link.

## Build

```sh
./build.sh
```

It runs the Go and webapp tests and leaves `dist/com.github.adriadam10.music-cards-<version>.tar.gz`.
Only Docker is needed locally. CI builds every push; pushing a tag `vX.Y.Z` equal to
the `version` in `plugin.json` publishes the release.

## Install

System Console → Plugins → Plugin Management → Upload Plugin, or:

```sh
mmctl plugin install-url https://github.com/adriadam10/mattermost-plugin-music-cards/releases/download/vX.Y.Z/com.github.adriadam10.music-cards-X.Y.Z.tar.gz
mmctl plugin enable com.github.adriadam10.music-cards
```
