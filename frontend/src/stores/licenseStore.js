import { writable } from 'svelte/store';
import { clientApi } from '../lib/api/client';
import { licenseApi } from '../lib/api/license';
import { logApi } from '../lib/api/log';

function createLicenseStore() {
    const { subscribe, set, update } = writable({
        stats: {
            total_clients: 0,
            active_licenses: 0,
            suspended_licenses: 0,
            unassigned_licenses: 0,
            revoked_licenses: 0
        },
        clients: [],
        licenses: [],
        logs: [],
        isLoading: false
    });

    async function fetchAll() {
        update(state => ({ ...state, isLoading: true }));
        try {
            const [statsRes, licensesRes, clientsRes, logsRes] = await Promise.all([
                clientApi.getStats(),
                licenseApi.getAll(),
                clientApi.getAll(),
                logApi.getAll()
            ]);

            set({
                stats: statsRes.data,
                licenses: licensesRes.data,
                clients: clientsRes.data,
                logs: logsRes.data,
                isLoading: false
            });
        } catch (err) {
            update(state => ({ ...state, isLoading: false }));
            throw err;
        }
    }

    return {
        subscribe,
        fetchAll,
        // Client operations
        createClient: async (name, owner, phone) => {
            await clientApi.create(name, owner, phone);
            await fetchAll();
        },
        updateClient: async (id, name, owner, phone) => {
            await clientApi.update(id, name, owner, phone);
            await fetchAll();
        },
        deleteClient: async (id) => {
            await clientApi.delete(id);
            await fetchAll();
        },
        // License operations
        generateLicense: async (clientId, trialLimit, expiresAt) => {
            const res = await licenseApi.generate(clientId, trialLimit, expiresAt);
            await fetchAll();
            return res.data;
        },
        updateLicenseStatus: async (id, status) => {
            await licenseApi.updateStatus(id, status);
            await fetchAll();
        },
        deleteLicense: async (id) => {
            await licenseApi.delete(id);
            await fetchAll();
        }
    };
}

export const licenseStore = createLicenseStore();
