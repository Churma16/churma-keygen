import { writable } from 'svelte/store';

function createToastStore() {
    const { subscribe, set } = writable({ message: '', type: 'success' });
    let timeoutId;

    return {
        subscribe,
        show: (message, type = 'success') => {
            set({ message, type });
            if (timeoutId) clearTimeout(timeoutId);
            timeoutId = setTimeout(() => {
                set({ message: '', type: 'success' });
            }, 3500);
        },
        clear: () => {
            set({ message: '', type: 'success' });
        }
    };
}

export const toast = createToastStore();
