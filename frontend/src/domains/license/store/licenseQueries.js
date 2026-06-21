import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { licenseApi } from '../api/licenseApi';

export function useLicensesQuery() {
    return createQuery({
        queryKey: ['licenses'],
        queryFn: async () => {
            const res = await licenseApi.getAll();
            return res.data || [];
        }
    });
}

export function useGenerateLicenseMutation() {
    const queryClient = useQueryClient();
    return createMutation({
        mutationFn: async ({ clientId, trialLimit, expiresAt }) => {
            const res = await licenseApi.generate(clientId, trialLimit, expiresAt);
            return res.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['licenses'] });
            queryClient.invalidateQueries({ queryKey: ['clientStats'] });
        }
    });
}

export function useUpdateLicenseStatusMutation() {
    const queryClient = useQueryClient();
    return createMutation({
        mutationFn: async ({ id, status }) => {
            return await licenseApi.updateStatus(id, status);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['licenses'] });
            queryClient.invalidateQueries({ queryKey: ['clientStats'] });
        }
    });
}

export function useDeleteLicenseMutation() {
    const queryClient = useQueryClient();
    return createMutation({
        mutationFn: async (id) => {
            return await licenseApi.delete(id);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['licenses'] });
            queryClient.invalidateQueries({ queryKey: ['clientStats'] });
        }
    });
}
