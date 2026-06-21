<script>
    import { formatDate } from '../../../shared/utils/format';

    export let filteredLogs = [];
</script>

<div class="overflow-x-auto w-full">
    <table class="table table-zebra w-full text-left">
        <thead>
        <tr class="bg-base-100/50 text-primary border-b border-base-300">
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">IP Address</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Kunci Lisensi</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Klien Terkait</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Hardware ID Percobaan</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Status Respon</th>
            <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">Waktu Upaya</th>
        </tr>
        </thead>
        <tbody class="divide-y divide-base-300">
        {#if filteredLogs.length === 0}
            <tr>
                <td colspan="6" class="text-center py-16 text-gray-400 font-medium">Tidak ada log audit yang terekam.</td>
            </tr>
        {:else}
            {#each filteredLogs as l}
                <tr class="hover:bg-base-100/40 transition text-xs">
                    <td class="py-3 pl-6 font-mono font-bold text-gray-700">{l.ip_address}</td>
                    <td class="py-3 font-mono font-bold text-primary">{l.attempted_code}</td>
                    <td class="py-3 font-semibold text-gray-800">
                        {#if l.license && l.license.client}
                            {l.license.client.name}
                        {:else}
                            <span class="text-gray-400 font-normal italic">Tidak terikat</span>
                        {/if}
                    </td>
                    <td class="py-3 font-mono text-[10px] max-w-[200px] truncate text-gray-500" title={l.hardware_id_attempt}>
                        {l.hardware_id_attempt}
                    </td>
                    <td class="py-3">
                        <span class="badge badge-sm font-bold text-[9px] px-1.5 py-0.5 border-none
                            {l.status === 'SUCCESS' ? 'text-primary bg-success/50' : ''}
                            {l.status === 'INVALID_KEY' ? 'text-red-800 bg-error/70' : ''}
                            {l.status === 'HWID_MISMATCH' ? 'text-yellow-800 bg-yellow-100' : ''}
                            {l.status === 'SUSPENDED_KEY' ? 'text-amber-800 bg-amber-100' : ''}
                        ">
                            {l.status}
                        </span>
                    </td>
                    <td class="py-3 text-right pr-6 text-gray-400 font-semibold">{formatDate(l.created_at)}</td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
</div>
