import React from 'react';

import {MusicIcon} from './icons.jsx';
import {safeCover} from './track.js';

// Album art, or a note when there is none (or the URL is not YouTube's).
export function Cover({url, size}) {
    const safe = safeCover(url);
    return safe ? (
        <img
            className='mc-cover-img'
            src={safe}
            width={size}
            height={size}
            alt=''
            loading='lazy'
            referrerPolicy='no-referrer'
        />
    ) : (
        <span
            className='mc-cover-img mc-cover-img--empty'
            style={{width: size, height: size}}
        >
            <MusicIcon size={size / 2}/>
        </span>
    );
}
