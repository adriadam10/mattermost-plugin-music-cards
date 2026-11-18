// Helpers for the song in a card's post props. Props can be set by anyone who
// can call the posts API, so everything that becomes a URL is checked here.

const VIDEO_ID = /^[A-Za-z0-9_-]{11}$/;
const COVER_HOSTS = new Set(['yt3.googleusercontent.com', 'lh3.googleusercontent.com', 'i.ytimg.com']);

export function trackFromProps(props) {
    const videoId = String(props?.video_id || '');
    if (!VIDEO_ID.test(videoId)) {
        return null;
    }
    const text = (key) => (typeof props[key] === 'string' ? props[key] : '');
    return {
        videoId,
        title: text('title'),
        artists: text('artists'),
        album: text('album'),
        duration: text('duration'),
        cover: safeCover(text('cover')),
        caption: text('caption'),
    };
}

export function safeCover(url) {
    try {
        const parsed = new URL(url);
        return parsed.protocol === 'https:' && COVER_HOSTS.has(parsed.hostname) ? parsed.href : '';
    } catch {
        return '';
    }
}

export const musicUrl = (videoId) => `https://music.youtube.com/watch?v=${videoId}`;

// youtube-nocookie.com does not set tracking cookies until the video plays.
export const embedUrl = (videoId) => `https://www.youtube-nocookie.com/embed/${videoId}?autoplay=1&rel=0`;

// "MOTOMAMI + · 2:37"
export const details = (track) => [track.album, track.duration].filter(Boolean).join(' · ');
