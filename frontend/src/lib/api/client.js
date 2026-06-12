import { baseApi } from './base';

export const clientApi = {
    getAll: () => baseApi.get('/admin/clients'),
    getStats: () => baseApi.get('/admin/stats'),
    create: (name, ownerName, phone) => baseApi.post('/admin/clients', { name, owner_name: ownerName, phone }),
    update: (id, name, ownerName, phone) => baseApi.put(`/admin/clients/${id}`, { name, owner_name: ownerName, phone }),
    delete: (id) => baseApi.delete(`/admin/clients/${id}`),
};
