import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { clientApi } from '../api/clientApi';

export function useClientsQuery() {
    return createQuery(() => ({
        queryKey: ['clients'],
        queryFn: async () => {
            const res = await clientApi.getAll();
            return res.data || [];
        }
    }));
}

export function useClientStatsQuery() {
    return createQuery(() => ({
        queryKey: ['clientStats'],
        queryFn: async () => {
            const res = await clientApi.getStats();
            return res.data;
        }
    }));
}

export function useCreateClientMutation() {
    const queryClient = useQueryClient();
    return createMutation(() => ({
        mutationFn: async ({ name, ownerName, phone }) => {
            return await clientApi.create(name, ownerName, phone);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['clients'] });
            queryClient.invalidateQueries({ queryKey: ['clientStats'] });
        }
    }));
}

export function useUpdateClientMutation() {
    const queryClient = useQueryClient();
    return createMutation(() => ({
        mutationFn: async ({ id, name, ownerName, phone }) => {
            return await clientApi.update(id, name, ownerName, phone);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['clients'] });
        }
    }));
}

export function useDeleteClientMutation() {
    const queryClient = useQueryClient();
    return createMutation(() => ({
        mutationFn: async (id) => {
            return await clientApi.delete(id);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['clients'] });
            queryClient.invalidateQueries({ queryKey: ['clientStats'] });
        }
    }));
}
