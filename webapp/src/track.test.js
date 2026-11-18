import assert from 'node:assert/strict';
import test from 'node:test';

import {details, safeCover, trackFromProps} from './track.js';

test('trackFromProps rejects a bad video id', () => {
    assert.equal(trackFromProps({video_id: 'x"onload=1'}), null);
    assert.equal(trackFromProps(undefined), null);
});

test('trackFromProps keeps only strings', () => {
    const track = trackFromProps({video_id: 'dN3y8kVNplA', title: 'DESPECHÁ', artists: {x: 1}});
    assert.equal(track.title, 'DESPECHÁ');
    assert.equal(track.artists, '');
});

test('safeCover only allows YouTube image hosts over https', () => {
    assert.equal(safeCover('https://yt3.googleusercontent.com/a=w544-h544'), 'https://yt3.googleusercontent.com/a=w544-h544');
    assert.equal(safeCover('https://i.ytimg.com/vi/x/hqdefault.jpg'), 'https://i.ytimg.com/vi/x/hqdefault.jpg');
    assert.equal(safeCover('http://i.ytimg.com/vi/x/hqdefault.jpg'), '');
    assert.equal(safeCover('https://evil.example/pixel.gif'), '');
    assert.equal(safeCover('javascript:alert(1)'), '');
    assert.equal(safeCover(''), '');
});

test('details', () => {
    assert.equal(details({album: 'LUX', duration: '3:16'}), 'LUX · 3:16');
    assert.equal(details({album: '', duration: '3:16'}), '3:16');
});
