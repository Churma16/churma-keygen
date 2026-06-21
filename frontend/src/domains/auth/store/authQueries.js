import { createMutation } from '@tanstack/svelte-query';
import { authApi } from '../api/authApi';

export function useUpdateProfileMutation() {
    return createMutation({
        mutationFn: async ({ username, current_password, new_password }) => {
            const res = await authApi.updateProfile({ username, current_password, new_password });
            return res.data;
        }
    });
}
