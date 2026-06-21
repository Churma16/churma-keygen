import { writable } from 'svelte/store';
import { logApi } from '../api/logApi';

function createLogStore() {
    const { subscribe, set, update } = writable({
        logs: [],
        isLoading: false
    });

    async function fetchLogs() {
        update(state => ({ ...state, isLoading: true }));
        try {
            const res = await logApi.getAll();
            update(state => ({
                logs: res.data,
                isLoading: false
            }));
        } catch (err) {
            update(state => ({ ...state, isLoading: false }));
            throw err;
        }
    }

    return {
        subscribe,
        fetchLogs
    };
}

export const logStore = createLogStore();
