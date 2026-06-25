import { baseApi } from '../../../shared/infra/api/baseApi';

export const settingApi = {
    get: (key) => baseApi.get(`/admin/settings/${key}`),
    update: (key, value) => baseApi.put(`/admin/settings/${key}`, { value })
};
