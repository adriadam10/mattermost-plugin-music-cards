import manifest from '../../plugin.json';

// Mattermost rejects cookie-authenticated requests without its CSRF header.
function csrfToken() {
    return document.cookie.match(/(?:^|;\s*)MMCSRF=([^;]+)/)?.[1] ?? '';
}

const base = () => `${window.basename || ''}/plugins/${manifest.id}/api/v1`;

async function request(method, path, body, signal) {
    const response = await fetch(`${base()}${path}`, {
        method,
        signal,
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest',
            'X-CSRF-Token': csrfToken(),
            ...(body ? {'Content-Type': 'application/json'} : {}),
        },
        body: body && JSON.stringify(body),
    });
    if (!response.ok) {
        throw new Error((await response.text()).trim() || `HTTP ${response.status}`);
    }
    return response.json();
}

export const searchSongs = (query, signal) => request('GET', `/search?${new URLSearchParams({q: query})}`, null, signal);

export const shareSong = ({channelId, rootId, videoId, text}) => request('POST', '/share', {
    channel_id: channelId,
    root_id: rootId || '',
    video_id: videoId,
    text: text || '',
});
