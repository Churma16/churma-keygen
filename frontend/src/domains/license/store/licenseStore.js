import { writable } from 'svelte/store';
import { licenseApi } from '../api/licenseApi';

function createLicenseStore() {
    const { subscribe, set, update } = writable({
        licenses: [],
        isLoading: false
    });

    async function fetchLicenses() {
        update(state => ({ ...state, isLoading: true }));
        try {
            const res = await licenseApi.getAll();
            update(state => ({
                licenses: res.data,
                isLoading: false
            }));
        } catch (err) {
            update(state => ({ ...state, isLoading: false }));
            throw err;
        }
    }

    return {
        subscribe,
        fetchLicenses,
        generateLicense: async (clientId, trialLimit, expiresAt) => {
            const res = await licenseApi.generate(clientId, trialLimit, expiresAt);
            await fetchLicenses();
            return res.data;
        },
        updateLicenseStatus: async (id, status) => {
            await licenseApi.updateStatus(id, status);
            await fetchLicenses();
        },
        deleteLicense: async (id) => {
            await licenseApi.delete(id);
            await fetchLicenses();
        }
    };
}

export const licenseStore = createLicenseStore();
