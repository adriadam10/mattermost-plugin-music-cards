import assert from 'node:assert/strict';
import test from 'node:test';

import {keysOf, LANGUAGES, translator} from './i18n.js';

test('every language has the same keys', () => {
    const [first, ...rest] = LANGUAGES;
    for (const lang of rest) {
        assert.deepEqual(keysOf(lang).sort(), keysOf(first).sort(), lang);
    }
});

test('unknown languages fall back to English', () => {
    assert.equal(translator('fr')('play'), 'Play');
    assert.equal(translator('es-ES')('play'), 'Reproducir');
});
