<script>
    import { createEventDispatcher } from 'svelte';
    import { UserPlus, Plus, Copy } from 'lucide-svelte';

    const dispatch = createEventDispatcher();
</script>

<div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm flex flex-col justify-between">
    <div>
        <h3 class="font-bold text-primary text-sm mb-1">Pintasan Aksi</h3>
        <p class="text-[11px] text-gray-400 font-medium mb-4">Gunakan pintasan cepat ini untuk mengelola lisensi klien dengan instan.</p>

        <div class="space-y-2">
            <button on:click={() => dispatch('addClient')}
                    class="btn btn-sm btn-outline btn-primary w-full gap-2 justify-start h-10 rounded-md text-xs font-bold">
                <UserPlus size={14}/>
                <span>Tambah Klien Baru</span>
            </button>
            <button on:click={() => dispatch('printLicense')}
                    class="btn btn-sm btn-primary w-full gap-2 justify-start h-10 rounded-md text-xs font-bold text-white">
                <Plus size={14}/>
                <span>Cetak Kunci Lisensi</span>
            </button>
            <button on:click={() => dispatch('copyPublicKey')}
                    class="btn btn-sm btn-outline btn-ghost w-full gap-2 justify-start h-10 rounded-md text-xs font-semibold text-gray-600 border-base-300 hover:bg-base-100">
                <Copy size={14}/>
                <span>Salin RSA Public Key</span>
            </button>
        </div>
    </div>

    <!-- Connection Status -->
    <div class="mt-6 border-t border-base-300 pt-4 flex items-center gap-3">
        <span class="relative flex h-2.5 w-2.5 shrink-0">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-success"></span>
        </span>
        <div class="text-[11px] font-semibold text-gray-500">
            Koneksi database PostgreSQL terhubung dan aman.
        </div>
    </div>
</div>
