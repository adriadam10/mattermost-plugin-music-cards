import React from 'react';

import {searchSongs, shareSong} from './api.js';
import {Cover} from './cover.jsx';
import {MusicIcon} from './icons.jsx';
import {details} from './track.js';

const DEBOUNCE_MS = 350;
const MIN_QUERY = 2;

function useSongSearch(query) {
    const [state, setState] = React.useState({tracks: [], loading: false, error: ''});
    React.useEffect(() => {
        const q = query.trim();
        if (q.length < MIN_QUERY) {
            setState({tracks: [], loading: false, error: ''});
            return undefined;
        }
        const abort = new AbortController();
        setState((s) => ({...s, loading: true, error: ''}));
        const timer = setTimeout(async () => {
            try {
                const {tracks} = await searchSongs(q, abort.signal);
                setState({tracks, loading: false, error: ''});
            } catch (err) {
                if (!abort.signal.aborted) {
                    setState({tracks: [], loading: false, error: err.message});
                }
            }
        }, DEBOUNCE_MS);
        return () => {
            clearTimeout(timer);
            abort.abort();
        };
    }, [query]);
    return state;
}

// Opened from the editor button: searches YouTube Music as you type and posts
// the chosen song, with whatever was written in the message box.
export default function Search({request, onDone, t}) {
    const [query, setQuery] = React.useState('');
    const [active, setActive] = React.useState(0);
    const [sending, setSending] = React.useState(false);
    const [sendError, setSendError] = React.useState('');
    const {tracks, loading, error} = useSongSearch(query);
    const caption = request.text?.trim() || '';

    React.useEffect(() => setActive(0), [tracks]);
    React.useEffect(() => {
        const onKey = (event) => event.key === 'Escape' && !sending && onDone();
        document.addEventListener('keydown', onKey);
        return () => document.removeEventListener('keydown', onKey);
    }, [onDone, sending]);

    const send = async (track) => {
        setSending(true);
        setSendError('');
        try {
            await shareSong({channelId: request.channelId, rootId: request.rootId, videoId: track.video_id, text: caption});
            if (caption) {
                request.clearText();
            }
            onDone();
        } catch (err) {
            setSending(false);
            setSendError(`${t('sendError')} ${err.message}`);
        }
    };

    const onKeyDown = (event) => {
        if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
            event.preventDefault();
            const step = event.key === 'ArrowDown' ? 1 : -1;
            setActive((i) => (i + step + tracks.length) % Math.max(tracks.length, 1));
        } else if (event.key === 'Enter' && tracks[active] && !sending) {
            event.preventDefault();
            send(tracks[active]);
        }
    };

    let status = '';
    if (sending) {
        status = t('sending');
    } else if (sendError || error) {
        status = sendError || `${t('searchError')} ${error}`;
    } else if (loading) {
        status = t('searching');
    } else if (query.trim().length >= MIN_QUERY && !tracks.length) {
        status = t('noResults');
    }

    return (
        <div
            className='mc-backdrop'
            onMouseDown={(event) => event.target === event.currentTarget && !sending && onDone()}
        >
            <div
                className='mc-modal'
                role='dialog'
                aria-modal='true'
                aria-label={t('title')}
            >
                <div className='mc-modal__header'>
                    <MusicIcon size={20}/>
                    <h2 className='mc-modal__title'>{t('title')}</h2>
                    <button
                        type='button'
                        className='mc-modal__close'
                        aria-label={t('close')}
                        onClick={onDone}
                        disabled={sending}
                    >
                        {'×'}
                    </button>
                </div>
                <div className='mc-modal__body'>
                    <input
                        className='mc-input'
                        type='search'
                        value={query}
                        placeholder={t('placeholder')}
                        autoFocus={true}
                        disabled={sending}
                        onChange={(event) => setQuery(event.target.value)}
                        onKeyDown={onKeyDown}
                    />
                    {caption && (
                        <p className='mc-note'>
                            {t('withText')} <span className='mc-quote'>{caption}</span>
                        </p>
                    )}
                    {status && <p className={sendError || error ? 'mc-error' : 'mc-note'}>{status}</p>}
                    <ul
                        className='mc-results'
                        role='listbox'
                    >
                        {tracks.map((track, i) => (
                            <li key={track.video_id}>
                                <button
                                    type='button'
                                    role='option'
                                    aria-selected={i === active}
                                    className={`mc-result${i === active ? ' mc-result--active' : ''}`}
                                    disabled={sending}
                                    onMouseEnter={() => setActive(i)}
                                    onClick={() => send(track)}
                                >
                                    <Cover
                                        url={track.cover}
                                        size={44}
                                    />
                                    <span className='mc-result__text'>
                                        <span className='mc-ellipsis mc-strong'>{track.title}</span>
                                        <span className='mc-ellipsis mc-muted'>{track.artists}</span>
                                    </span>
                                    <span className='mc-ellipsis mc-muted mc-result__details'>{details(track)}</span>
                                </button>
                            </li>
                        ))}
                    </ul>
                    {tracks.length > 0 && <p className='mc-note mc-hint'>{t('hint')}</p>}
                </div>
            </div>
        </div>
    );
}
