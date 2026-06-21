import { baseApi } from '../../../shared/infra/api/baseApi';

export const logApi = {
    getAll: () => baseApi.get('/admin/logs'),
};
