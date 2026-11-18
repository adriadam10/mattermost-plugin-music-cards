// Mattermost's theme variables, so the plugin follows the user's theme.
const MUTED = 'rgba(var(--center-channel-color-rgb), 0.64)';
const BORDER = 'rgba(var(--center-channel-color-rgb), 0.16)';
const HOVER = 'rgba(var(--center-channel-color-rgb), 0.08)';

export const STYLES = `
.mc-editor-btn {
    display: inline-flex;
    width: 32px;
    height: 32px;
    align-items: center;
    justify-content: center;
    padding: 0;
    border: none;
    border-radius: 4px;
    background: none;
    color: ${MUTED};
    cursor: pointer;
}
.mc-editor-btn:hover { background: ${HOVER}; color: var(--center-channel-color); }
.mc-backdrop {
    position: fixed;
    inset: 0;
    z-index: 10000;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding-top: 10vh;
    background: rgba(0, 0, 0, 0.5);
}
.mc-modal {
    display: flex;
    width: min(560px, calc(100vw - 32px));
    max-height: 80vh;
    flex-direction: column;
    border-radius: 8px;
    background: var(--center-channel-bg);
    box-shadow: 0 20px 32px rgba(0, 0, 0, 0.12);
    color: var(--center-channel-color);
}
.mc-modal__header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 16px 20px;
    border-bottom: 1px solid ${BORDER};
}
.mc-modal__title { flex: 1; margin: 0; font-size: 18px; font-weight: 600; }
.mc-modal__close {
    padding: 0 4px;
    border: none;
    background: none;
    color: ${MUTED};
    font-size: 24px;
    line-height: 1;
    cursor: pointer;
}
.mc-modal__body {
    display: flex;
    min-height: 0;
    flex-direction: column;
    gap: 10px;
    padding: 16px 20px;
}
.mc-input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid ${BORDER};
    border-radius: 4px;
    background: var(--center-channel-bg);
    color: var(--center-channel-color);
    font: inherit;
}
.mc-input:focus { border-color: var(--button-bg); outline: none; box-shadow: inset 0 0 0 1px var(--button-bg); }
.mc-note { margin: 0; color: ${MUTED}; font-size: 13px; }
.mc-hint { text-align: right; font-size: 12px; }
.mc-quote { color: var(--center-channel-color); font-style: italic; }
.mc-error { margin: 0; color: var(--error-text, #d24b4e); font-size: 13px; }
.mc-results {
    display: flex;
    min-height: 0;
    flex-direction: column;
    gap: 2px;
    margin: 0 -8px;
    padding: 0;
    overflow-y: auto;
    list-style: none;
}
.mc-result {
    display: flex;
    width: 100%;
    align-items: center;
    gap: 12px;
    padding: 6px 8px;
    border: none;
    border-radius: 6px;
    background: none;
    color: inherit;
    font: inherit;
    text-align: left;
    cursor: pointer;
}
.mc-result--active { background: ${HOVER}; }
.mc-result:disabled { opacity: 0.6; cursor: default; }
.mc-result__text { display: flex; min-width: 0; flex: 1; flex-direction: column; }
.mc-result__details { max-width: 40%; font-size: 12px; }
.mc-ellipsis { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mc-strong { font-weight: 600; }
.mc-muted { color: ${MUTED}; }
.mc-cover-img { display: block; flex-shrink: 0; border-radius: 4px; object-fit: cover; }
.mc-cover-img--empty {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: ${HOVER};
    color: ${MUTED};
}
.mc-post { display: flex; flex-direction: column; gap: 6px; }
.mc-caption--plain { white-space: pre-wrap; overflow-wrap: anywhere; }
.mc-card {
    display: flex;
    width: 420px;
    max-width: 100%;
    flex-direction: column;
    gap: 10px;
    margin: 2px 0;
    padding: 10px;
    border: 1px solid ${BORDER};
    border-radius: 8px;
    background: rgba(var(--center-channel-color-rgb), 0.03);
}
.mc-card__row { display: flex; align-items: center; gap: 12px; }
.mc-card__cover {
    position: relative;
    flex-shrink: 0;
    padding: 0;
    border: none;
    border-radius: 4px;
    background: none;
    cursor: pointer;
    overflow: hidden;
}
.mc-card__play {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.35);
    color: #fff;
    opacity: 0;
    transition: opacity 0.15s;
}
.mc-card__cover:hover .mc-card__play,
.mc-card__cover:focus-visible .mc-card__play,
.mc-card__cover--playing .mc-card__play { opacity: 1; }
.mc-card__info { display: flex; min-width: 0; flex: 1; flex-direction: column; gap: 2px; font-size: 14px; }
.mc-card__info .mc-muted { font-size: 12px; }
.mc-card__title, .mc-card__title:hover { color: var(--center-channel-color); text-decoration: none; }
.mc-card__title:hover { text-decoration: underline; }
.mc-card__open {
    display: inline-flex;
    flex-shrink: 0;
    align-self: flex-start;
    padding: 6px;
    border-radius: 4px;
    color: ${MUTED};
}
.mc-card__open:hover { background: ${HOVER}; color: var(--center-channel-color); }
.mc-card__player { position: relative; aspect-ratio: 16 / 9; border-radius: 6px; overflow: hidden; background: #000; }
.mc-card__player iframe { position: absolute; inset: 0; width: 100%; height: 100%; border: 0; }
.mc-eq { display: flex; height: 18px; align-items: flex-end; gap: 3px; }
.mc-eq i { width: 4px; border-radius: 1px; background: currentColor; animation: mc-eq 0.9s ease-in-out infinite; }
.mc-eq i:nth-child(2) { animation-delay: -0.3s; }
.mc-eq i:nth-child(3) { animation-delay: -0.6s; }
@keyframes mc-eq { 0%, 100% { height: 30%; } 50% { height: 100%; } }
@media (prefers-reduced-motion: reduce) { .mc-eq i { animation: none; height: 70%; } }
`;
