<script>
    import { push } from 'svelte-spa-router';
    
    // TanStack Query Hooks
    import { useCreateClientMutation } from '../../client/store/clientQueries';
    import { useLogsQuery } from '../../log/store/logQueries';
    
    // Shared Utils & State
    import { formatDate } from '../../../shared/utils/format';
    import { toast } from '../../../shared/store/toastStore';
    import { licenseApi } from '../../license/api/licenseApi';

    // Domain Components
    import StatsGrid from '../../license/components/StatsGrid.svelte';
    import QuickActions from '../../license/components/QuickActions.svelte';
    import ModalGenerator from '../../license/components/ModalGenerator.svelte';
    import CreateClientModal from '../../client/components/CreateClientModal.svelte';

    // Modals state
    let showGenModal = $state(false);
    let selectedClientForLicense = $state(null);
    let showCreateClientModal = $state(false);

    // Initialize TanStack Query Hooks
    const logsQuery = useLogsQuery();
    const createClientMutation = useCreateClientMutation();

    function openGenModalDirect(client = null) {
        selectedClientForLicense = client;
        showGenModal = true;
    }

    async function copyPublicKey() {
        try {
            const text = await licenseApi.getPublicKey();
            await navigator.clipboard.writeText(text);
            toast.show('Public Key RSA disalin! Siap dipasang di biner POS klien.', 'info');
        } catch (e) {
            toast.show('Gagal mengunduh Public Key.', 'error');
        }
    }

    async function handleCreateClient(event) {
        const { name, ownerName, phone } = event.detail;
        try {
            await createClientMutation.mutateAsync({ name, ownerName, phone });
            toast.show('Klien baru berhasil ditambahkan.');
            showCreateClientModal = false;
        } catch (e) {
            toast.show(e.message || 'Gagal menambahkan klien.', 'error');
        }
    }
</script>

<div class="mb-8">
    <h2 class="text-2xl font-bold text-primary tracking-tight">Selamat Datang, Admin</h2>
    <p class="text-xs text-gray-500 mt-1 font-medium">Berikut adalah ringkasan performa lisensi dan audit sistem hari ini.</p>
</div>

<!-- STATS GRID COMPONENT -->
<StatsGrid />

<!-- Overview Dashboard Widgets -->
<div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
    <!-- QUICK ACTIONS COMPONENT -->
    <QuickActions 
        on:addClient={() => showCreateClientModal = true}
        on:printLicense={() => openGenModalDirect()}
        on:copyPublicKey={copyPublicKey}
    />

    <!-- Right: Recent Activity Log summary -->
    <div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm lg:col-span-2">
        <div class="flex items-center justify-between mb-4">
            <div>
                <h3 class="font-bold text-primary text-sm mb-1">Aktivitas Aktivasi Terkini</h3>
                <p class="text-[11px] text-gray-400 font-medium">Log upaya pencocokan serial kunci dan HWID oleh komputer kasir lokal.</p>
            </div>
            <button on:click={() => push('/logs')}
                    class="btn btn-xs btn-ghost text-xs text-primary font-bold hover:bg-base-100">
                Lihat Semua
            </button>
        </div>

        <div class="overflow-x-auto">
            <table class="table table-compact w-full text-xs">
                <thead>
                <tr class="text-gray-400 border-b border-base-300">
                    <th class="py-2.5 font-bold">Klien</th>
                    <th class="py-2.5 font-bold">IP</th>
                    <th class="py-2.5 font-bold">Status</th>
                    <th class="py-2.5 font-bold text-right">Waktu</th>
                </tr>
                </thead>
                <tbody class="divide-y divide-base-300/65">
                {#if (logsQuery.data || []).length === 0}
                    <tr>
                        <td colspan="4" class="text-center py-8 text-gray-400 font-medium">Tidak ada rekaman log.</td>
                    </tr>
                {:else}
                    {#each (logsQuery.data || []).slice(0, 5) as l}
                        <tr>
                            <td class="py-2.5 font-semibold text-primary">
                                {#if l.license && l.license.client}
                                    {l.license.client.name}
                                {:else}
                                    <span class="text-gray-400 font-normal italic">Kunci salah</span>
                                {/if}
                            </td>
                            <td class="py-2.5 font-mono">{l.ip_address}</td>
                            <td class="py-2.5">
                                <span class="badge badge-sm font-bold text-[8px] px-1.5 py-0.5 border-none
                                    {l.status === 'SUCCESS' ? 'badge-success text-primary bg-success/40' : ''}
                                    {l.status === 'INVALID_KEY' ? 'badge-error text-red-700 bg-error/70' : ''}
                                    {l.status === 'HWID_MISMATCH' ? 'badge-warning text-yellow-800 bg-yellow-100' : ''}
                                    {l.status === 'SUSPENDED_KEY' ? 'badge-neutral text-amber-700 bg-amber-100' : ''}
                                ">
                                    {l.status}
                                </span>
                            </td>
                            <td class="py-2.5 text-right text-gray-400 font-medium">{formatDate(l.created_at)}</td>
                        </tr>
                    {/each}
                {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>

<!-- Modal: Create Client -->
{#if showCreateClientModal}
    <CreateClientModal 
        isSubmitting={createClientMutation.isPending}
        on:close={() => showCreateClientModal = false}
        on:submit={handleCreateClient}
    />
{/if}

<!-- License Generator Modal integration -->
{#if showGenModal}
    <ModalGenerator
            selectedClient={selectedClientForLicense}
            on:close={() => showGenModal = false}
    />
{/if}
