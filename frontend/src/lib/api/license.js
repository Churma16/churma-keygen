import { baseApi } from './base';

export const licenseApi = {
    getAll: () => baseApi.get('/admin/licenses'),
    generate: (clientId, trialLimit, expiresAt) => baseApi.post('/admin/licenses', {
        client_id: clientId,
        trial_limit: trialLimit,
        expires_at: expiresAt
    }),
    updateStatus: (id, status) => baseApi.put(`/admin/licenses/${id}/status`, { status }),
    delete: (id) => baseApi.delete(`/admin/licenses/${id}`),
    getPublicKey: () => baseApi.get('/public-key'),
};
