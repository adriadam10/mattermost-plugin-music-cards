import React from 'react';

import {OpenIcon, PlayIcon} from './icons.jsx';
import {Cover} from './cover.jsx';
import {details, embedUrl, musicUrl, trackFromProps} from './track.js';

// The user's text, with Mattermost's own markdown, mentions and emoji.
function Caption({text}) {
    const utils = window.PostUtils;
    if (!utils) {
        return <div className='mc-caption mc-caption--plain'>{text}</div>;
    }
    return <div className='mc-caption'>{utils.messageHtmlToComponent(utils.formatText(text, {}), {})}</div>;
}

// Renders a custom_music_card post. The player only loads (and only reaches
// YouTube) when the play button is pressed; starting another card stops it.
export default function Card({post, bus, t}) {
    const track = trackFromProps(post.props);
    const [playing, setPlaying] = React.useState(false);

    React.useEffect(() => bus.on('play', (id) => id !== post.id && setPlaying(false)), [bus, post.id]);

    if (!track) {
        return <div className='mc-caption mc-caption--plain'>{post.message}</div>;
    }

    const toggle = () => {
        if (!playing) {
            bus.emit('play', post.id);
        }
        setPlaying(!playing);
    };
    const link = musicUrl(track.videoId);

    return (
        <div className='mc-post'>
            {track.caption && <Caption text={track.caption}/>}
            <div className='mc-card'>
                <div className='mc-card__row'>
                    <button
                        type='button'
                        className={`mc-card__cover${playing ? ' mc-card__cover--playing' : ''}`}
                        title={playing ? t('stop') : t('play')}
                        aria-label={playing ? t('stop') : t('play')}
                        onClick={toggle}
                    >
                        <Cover
                            url={track.cover}
                            size={72}
                        />
                        <span className='mc-card__play'>
                            {playing ? <span className='mc-eq'><i/><i/><i/></span> : <PlayIcon/>}
                        </span>
                    </button>
                    <div className='mc-card__info'>
                        <a
                            className='mc-ellipsis mc-strong mc-card__title'
                            href={link}
                            target='_blank'
                            rel='noopener noreferrer'
                        >
                            {track.title}
                        </a>
                        <span className='mc-ellipsis'>{track.artists}</span>
                        <span className='mc-ellipsis mc-muted'>{details(track)}</span>
                    </div>
                    <a
                        className='mc-card__open'
                        href={link}
                        target='_blank'
                        rel='noopener noreferrer'
                        title={t('openInMusic')}
                        aria-label={t('openInMusic')}
                    >
                        <OpenIcon/>
                    </a>
                </div>
                {playing && (
                    <div className='mc-card__player'>
                        <iframe
                            src={embedUrl(track.videoId)}
                            title={track.title}
                            allow='autoplay; encrypted-media; picture-in-picture'
                            allowFullScreen={true}
                            referrerPolicy='strict-origin-when-cross-origin'
                        />
                    </div>
                )}
            </div>
        </div>
    );
}
