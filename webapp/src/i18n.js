const STRINGS = {
    en: {
        button: 'Share a song',
        title: 'Share a song',
        placeholder: 'Search a song or paste a YouTube Music link',
        withText: 'It is sent with the text in your message box:',
        searching: 'Searching…',
        noResults: 'No songs found.',
        searchError: 'Could not search:',
        sending: 'Sending…',
        sendError: 'Could not send it:',
        hint: '↑ ↓ to choose, Enter to send',
        close: 'Close',
        play: 'Play',
        stop: 'Close player',
        openInMusic: 'Open in YouTube Music',
    },
    es: {
        button: 'Compartir canción',
        title: 'Compartir canción',
        placeholder: 'Busca una canción o pega un enlace de YouTube Music',
        withText: 'Se envía con el texto que tienes escrito:',
        searching: 'Buscando…',
        noResults: 'No hay canciones.',
        searchError: 'No se ha podido buscar:',
        sending: 'Enviando…',
        sendError: 'No se ha podido enviar:',
        hint: '↑ ↓ para elegir, Intro para enviar',
        close: 'Cerrar',
        play: 'Reproducir',
        stop: 'Cerrar reproductor',
        openInMusic: 'Abrir en YouTube Music',
    },
};

export function translator(locale) {
    const strings = STRINGS[(locale || '').split('-')[0]] || STRINGS.en;
    return (key) => strings[key];
}

export const LANGUAGES = Object.keys(STRINGS);
export const keysOf = (lang) => Object.keys(STRINGS[lang]);
