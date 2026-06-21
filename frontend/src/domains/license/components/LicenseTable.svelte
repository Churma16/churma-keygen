<script>
    import { createEventDispatcher } from 'svelte';
    import { ShieldAlert, ShieldCheck, ShieldX, Trash2 } from 'lucide-svelte';
    import { formatDate } from '../../../shared/utils/format';

    const dispatch = createEventDispatcher();

    export let filteredLicenses = [];

    function handleUpdateStatus(id, newStatus) {
        dispatch('updateStatus', { id, status: newStatus });
    }

    function handleDelete(id) {
        dispatch('deleteLicense', id);
    }
</script>

<div class="overflow-x-auto w-full">
    <table class="table table-zebra w-full text-left">
        <thead>
        <tr class="bg-base-100/50 text-primary border-b border-base-300">
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">Klien Toko</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Kode Lisensi</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Hardware ID</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Batas Transaksi</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Status</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Masa Aktif</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">Aksi</th>
        </tr>
        </thead>
        <tbody class="divide-y divide-base-300">
        {#if filteredLicenses.length === 0}
            <tr>
                <td colspan="7" class="text-center py-16 text-gray-400 font-medium">Tidak ada lisensi terdaftar yang cocok.</td>
            </tr>
        {:else}
            {#each filteredLicenses as lic}
                <tr class="hover:bg-base-100/40 transition">
                    <td class="py-4 font-semibold text-primary pl-6 max-w-[200px] truncate">
                        {lic.client_name || (lic.client ? lic.client.name : 'Unknown Client')}
                    </td>
                    <td class="py-4 font-mono font-bold text-gray-800 text-sm select-all">
                        {lic.license_code}
                    </td>
                    <td class="py-4 font-mono text-xs max-w-[150px] truncate text-gray-500">
                        {#if lic.hardware_id}
                            <span class="tooltip cursor-help" data-tip={lic.hardware_id}>
                                {lic.hardware_id.slice(0, 16)}...
                            </span>
                        {:else}
                            <span class="text-gray-400 italic">Belum terkunci</span>
                        {/if}
                    </td>
                    <td class="py-4 text-xs font-semibold">
                        {#if lic.trial_limit === 0 || lic.trial_limit === -1}
                            <span class="badge badge-success text-[10px] font-bold py-1.5 px-2 border-none text-primary bg-success/60">Unlimited</span>
                        {:else}
                            <span>{lic.trial_limit} Transaksi</span>
                        {/if}
                    </td>
                    <td class="py-4 text-xs">
                        {#if lic.status === 'ACTIVE'}
                            <span class="badge badge-success text-[9px] font-extrabold px-2 py-1 border-none text-primary bg-success/50">AKTIF</span>
                        {:else}
                            <span class="badge text-[9px] font-extrabold px-2 py-1 border-none
                                {lic.status === 'SUSPENDED' ? 'text-amber-800 bg-amber-100' : ''}
                                {lic.status === 'REVOKED' ? 'text-red-800 bg-error/70' : ''}
                                {lic.status === 'UNASSIGNED' ? 'text-gray-500 bg-gray-100' : ''}
                            ">
                                {lic.status}
                            </span>
                        {/if}
                    </td>
                    <td class="py-4 text-xs text-gray-400 font-semibold">
                        {#if lic.expires_at}
                            {formatDate(lic.expires_at)}
                        {:else}
                            <span class="text-gray-400 font-normal">Permanen</span>
                        {/if}
                    </td>
                    <td class="py-4 text-right pr-6">
                        <div class="flex justify-end gap-2">
                            {#if lic.status === 'ACTIVE'}
                                <button
                                        on:click={() => handleUpdateStatus(lic.id, 'SUSPENDED')}
                                        class="btn btn-xs btn-outline btn-warning gap-1 tooltip"
                                        data-tip="Tangguhkan Lisensi"
                                >
                                    <ShieldAlert size={12}/>
                                    <span>Tangguhkan</span>
                                </button>
                            {:else if lic.status === 'SUSPENDED'}
                                <button
                                        on:click={() => handleUpdateStatus(lic.id, 'ACTIVE')}
                                        class="btn btn-xs btn-outline btn-primary gap-1 tooltip"
                                        data-tip="Aktifkan Ulang"
                                >
                                    <ShieldCheck size={12}/>
                                    <span>Aktifkan</span>
                                </button>
                            {/if}

                            {#if lic.status !== 'REVOKED' && lic.status !== 'UNASSIGNED'}
                                <button
                                        on:click={() => handleUpdateStatus(lic.id, 'REVOKED')}
                                        class="btn btn-xs btn-outline btn-error gap-1 tooltip"
                                        data-tip="Batalkan Lisensi"
                                >
                                    <ShieldX size={12}/>
                                    <span>Batalkan</span>
                                </button>
                            {/if}

                            <button
                                    on:click={() => handleDelete(lic.id)}
                                    class="w-6 h-6 border border-red-100 hover:bg-red-50 text-red-500 rounded-md flex items-center justify-center tooltip"
                                    data-tip="Hapus Lisensi"
                            >
                                <Trash2 size={12}/>
                            </button>
                        </div>
                    </td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
</div>
