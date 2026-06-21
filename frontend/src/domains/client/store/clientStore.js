import { writable } from 'svelte/store';
import { clientApi } from '../api/clientApi';

function createClientStore() {
    const { subscribe, set, update } = writable({
        clients: [],
        stats: {
            total_clients: 0,
            active_licenses: 0,
            suspended_licenses: 0,
            unassigned_licenses: 0,
            revoked_licenses: 0
        },
        isLoading: false
    });

    async function fetchClients() {
        update(s => ({ ...s, isLoading: true }));
        try {
            const res = await clientApi.getAll();
            update(s => ({ ...s, clients: res.data, isLoading: false }));
        } catch (err) {
            update(s => ({ ...s, isLoading: false }));
            throw err;
        }
    }

    async function fetchStats() {
        update(s => ({ ...s, isLoading: true }));
        try {
            const res = await clientApi.getStats();
            update(s => ({ ...s, stats: res.data, isLoading: false }));
        } catch (err) {
            update(s => ({ ...s, isLoading: false }));
            throw err;
        }
    }

    return {
        subscribe,
        fetchClients,
        fetchStats,
        createClient: async (name, owner, phone) => {
            await clientApi.create(name, owner, phone);
            await fetchClients();
            await fetchStats();
        },
        updateClient: async (id, name, owner, phone) => {
            await clientApi.update(id, name, owner, phone);
            await fetchClients();
        },
        deleteClient: async (id) => {
            await clientApi.delete(id);
            await fetchClients();
            await fetchStats();
        }
    };
}

export const clientStore = createClientStore();
