import { createQuery, createMutation } from '@tanstack/svelte-query';
import { settingApi } from '../api/settingApi';

export function useGetSetting(key) {
    return createQuery(() => ({
        queryKey: ['settings', key],
        queryFn: async () => {
            const res = await settingApi.get(key);
            return res.data;
        }
    }));
}

export function useUpdateSettingMutation() {
    return createMutation(() => ({
        mutationFn: async ({ key, value }) => {
            const res = await settingApi.update(key, value);
            return res.data;
        }
    }));
}
