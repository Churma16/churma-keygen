import { baseApi } from './base';

export const logApi = {
    getAll: () => baseApi.get('/admin/logs'),
};
