import { baseApi } from '../../../shared/infra/api/baseApi';

export const authApi = {
    login: (username, password) => baseApi.post('/auth/login', { username, password }),
    getMe: () => baseApi.get('/admin/me'),
};
