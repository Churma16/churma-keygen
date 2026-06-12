import { writable } from 'svelte/store';

function createAuthStore() {
    const { subscribe, set, update } = writable({
        token: localStorage.getItem('admin_token') || '',
        username: localStorage.getItem('admin_username') || '',
        role: localStorage.getItem('admin_role') || '',
        isLoggedIn: !!localStorage.getItem('admin_token')
    });

    return {
        subscribe,
        login: (token, username, role) => {
            localStorage.setItem('admin_token', token);
            localStorage.setItem('admin_username', username);
            localStorage.setItem('admin_role', role);
            set({ token, username, role, isLoggedIn: true });
        },
        logout: () => {
            localStorage.removeItem('admin_token');
            localStorage.removeItem('admin_username');
            localStorage.removeItem('admin_role');
            set({ token: '', username: '', role: '', isLoggedIn: false });
            window.location.reload();
        }
    };
}

export const authStore = createAuthStore();
