import { createQuery } from '@tanstack/svelte-query';
import { logApi } from '../api/logApi';

export function useLogsQuery() {
    return createQuery({
        queryKey: ['logs'],
        queryFn: async () => {
            const res = await logApi.getAll();
            return res.data || [];
        }
    });
}
