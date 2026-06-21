<script>
    import { createEventDispatcher } from 'svelte';
    import { Plus, Edit, Trash2 } from 'lucide-svelte';
    import { formatDate } from '../../../shared/utils/format';

    const dispatch = createEventDispatcher();

    export let filteredClients = [];
</script>

<div class="overflow-x-auto w-full">
    <table class="table table-zebra w-full text-left">
        <thead>
        <tr class="bg-base-100/50 text-primary border-b border-base-300">
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">Nama Klien / Perusahaan</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Nama Pemilik</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Nomor Telepon</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Jumlah Lisensi</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Terdaftar Pada</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">Aksi</th>
        </tr>
        </thead>
        <tbody class="divide-y divide-base-300">
        {#if filteredClients.length === 0}
            <tr>
                <td colspan="6" class="text-center py-16 text-gray-400 font-medium">Tidak ada klien toko terdaftar.</td>
            </tr>
        {:else}
            {#each filteredClients as c}
                <tr class="hover:bg-base-100/40 transition text-sm">
                    <td class="py-4 font-semibold text-primary pl-6 select-all">{c.name}</td>
                    <td class="py-4 font-medium text-gray-700">
                        {#if c.owner_name}
                            {c.owner_name}
                        {:else}
                            <span class="text-gray-400 italic">Kosong</span>
                        {/if}
                    </td>
                    <td class="py-4 select-all text-xs font-semibold text-gray-600">
                        {#if c.phone}
                            {c.phone}
                        {:else}
                            <span class="text-gray-400 italic">Kosong</span>
                        {/if}
                    </td>
                    <td class="py-4 text-xs">
                        <span class="badge badge-success text-[10px] font-bold py-1.5 px-2 border-none text-primary bg-success/60">
                            {c.licenses ? c.licenses.length : 0} kunci
                        </span>
                    </td>
                    <td class="py-4 text-xs text-gray-400 font-semibold">{formatDate(c.created_at)}</td>
                    <td class="py-4 text-right pr-6">
                        <div class="flex justify-end gap-2">
                            <button
                                    on:click={() => dispatch('openGenModal', c)}
                                    class="btn btn-xs btn-primary gap-1 rounded-md text-white"
                            >
                                <Plus size={12}/>
                                <span>Cetak Kunci</span>
                            </button>

                            <button
                                    on:click={() => dispatch('editClient', c)}
                                    class="w-6 h-6 border border-base-300 hover:bg-base-100 text-gray-600 rounded-md flex items-center justify-center tooltip"
                                    data-tip="Edit Profil"
                            >
                                <Edit size={12}/>
                            </button>

                            <button
                                    on:click={() => dispatch('deleteClient', c.id)}
                                    class="w-6 h-6 border border-red-100 hover:bg-red-50 text-red-500 rounded-md flex items-center justify-center tooltip"
                                    data-tip="Hapus Klien"
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
