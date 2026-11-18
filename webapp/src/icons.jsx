import React from 'react';

const svgProps = (size) => ({
    viewBox: '0 0 24 24',
    width: size,
    height: size,
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: 2,
    strokeLinecap: 'round',
    strokeLinejoin: 'round',
    'aria-hidden': 'true',
});

export function MusicIcon({size = 18}) {
    return (
        <svg {...svgProps(size)}>
            <path d='M9 18V5l12-2v13'/>
            <circle
                cx='6'
                cy='18'
                r='3'
            />
            <circle
                cx='18'
                cy='16'
                r='3'
            />
        </svg>
    );
}

export function PlayIcon({size = 22}) {
    return (
        <svg
            viewBox='0 0 24 24'
            width={size}
            height={size}
            fill='currentColor'
            aria-hidden='true'
        >
            <path d='M8 5.5v13l11-6.5z'/>
        </svg>
    );
}

export function OpenIcon({size = 16}) {
    return (
        <svg {...svgProps(size)}>
            <path d='M14 4h6v6M20 4l-9 9M18 14v5a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1h5'/>
        </svg>
    );
}
