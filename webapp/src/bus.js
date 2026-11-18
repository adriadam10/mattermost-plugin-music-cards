// Lets the editor button open the search, which lives in the one root
// component, and makes sure only one card plays at a time.
export function createBus() {
    const listeners = new Map();
    return {
        on(event, listener) {
            if (!listeners.has(event)) {
                listeners.set(event, new Set());
            }
            listeners.get(event).add(listener);
            return () => listeners.get(event).delete(listener);
        },
        emit(event, payload) {
            listeners.get(event)?.forEach((listener) => listener(payload));
        },
    };
}
