import React from 'react';

import manifest from '../../plugin.json';

import {createBus} from './bus.js';
import Card from './card.jsx';
import {translator} from './i18n.js';
import {MusicIcon} from './icons.jsx';
import Search from './search.jsx';
import {STYLES} from './styles.js';

const POST_TYPE = 'custom_music_card';

function currentLocale(state) {
    const {currentUserId, profiles} = state.entities.users;
    return profiles[currentUserId]?.locale || document.documentElement.lang || 'en';
}

// The search window, opened from the editor button.
function Root({bus, t}) {
    const [request, setRequest] = React.useState(null);
    React.useEffect(() => bus.on('search', setRequest), [bus]);
    const close = React.useCallback(() => setRequest(null), []);
    return request && (
        <Search
            request={request}
            onDone={close}
            t={t}
        />
    );
}

// Rendered by Mattermost in the message box toolbar, with the box's draft.
function EditorButton({draft, updateText, bus, t}) {
    if (!draft?.channelId) {
        return null;
    }
    return (
        <button
            type='button'
            className='mc-editor-btn'
            title={t('button')}
            aria-label={t('button')}
            onClick={() => bus.emit('search', {
                channelId: draft.channelId,
                rootId: draft.rootId,
                text: draft.message,
                clearText: () => updateText(''),
            })}
        >
            <MusicIcon/>
        </button>
    );
}

class Plugin {
    initialize(registry, store) {
        const t = translator(currentLocale(store.getState()));
        const bus = createBus();

        this.style = document.createElement('style');
        this.style.textContent = STYLES;
        document.head.appendChild(this.style);

        registry.registerRootComponent(() => (
            <Root
                bus={bus}
                t={t}
            />
        ));
        registry.registerPostEditorActionComponent((props) => (
            <EditorButton
                {...props}
                bus={bus}
                t={t}
            />
        ));
        registry.registerPostTypeComponent(POST_TYPE, (props) => (
            <Card
                post={props.post}
                bus={bus}
                t={t}
            />
        ));
    }

    uninitialize() {
        this.style?.remove();
    }
}

window.registerPlugin(manifest.id, new Plugin());
