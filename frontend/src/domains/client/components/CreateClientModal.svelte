<script>
    import { createEventDispatcher } from 'svelte';
    import { X, UserPlus } from 'lucide-svelte';

    const dispatch = createEventDispatcher();

    export let isSubmitting = false;

    let name = '';
    let ownerName = '';
    let phone = '';

    function handleSave() {
        if (!name) return;
        dispatch('submit', { name, ownerName, phone });
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<div class="modal modal-open" on:click|self={() => dispatch('close')}>
    <div class="modal-box bg-white border border-base-300 rounded-lg shadow-2xl relative p-6 max-w-md">
        <button on:click={() => dispatch('close')}
                class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4 text-gray-400 hover:text-gray-600">
            <X size={16}/>
        </button>

        <div class="flex items-center gap-3 mb-5">
            <div class="p-2.5 bg-success/80 text-primary rounded-md border border-success">
                <UserPlus size={18}/>
            </div>
            <h3 class="font-bold text-lg text-primary">Daftarkan Klien Baru</h3>
        </div>

        <div class="space-y-4">
            <div class="form-control">
                <label class="label px-0.5 py-1" for="new-name">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Toko / Perusahaan</span>
                </label>
                <input
                        id="new-name"
                        type="text"
                        placeholder="Contoh: Samudera Furniture Padang"
                        bind:value={name}
                        class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                />
            </div>

            <div class="form-control">
                <label class="label px-0.5 py-1" for="new-owner">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Pemilik Toko</span>
                </label>
                <input
                        id="new-owner"
                        type="text"
                        placeholder="Contoh: Bpk. Muhammad"
                        bind:value={ownerName}
                        class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                />
            </div>

            <div class="form-control">
                <label class="label px-0.5 py-1" for="new-phone">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nomor Telepon</span>
                </label>
                <input
                        id="new-phone"
                        type="text"
                        placeholder="Contoh: 08123456789"
                        bind:value={phone}
                        class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                />
            </div>
        </div>

        <div class="modal-action gap-2 mt-6">
            <button on:click={() => dispatch('close')}
                    class="btn btn-outline btn-sm rounded-md text-xs font-semibold h-9 px-4">Batal
            </button>
            <button
                    on:click={handleSave}
                    disabled={!name || isSubmitting}
                    class="btn btn-primary btn-sm text-white rounded-md text-xs font-bold h-9 px-4"
            >
                {#if isSubmitting}
                    <span class="loading loading-spinner loading-xs mr-1"></span>
                {/if}
                Simpan Klien
            </button>
        </div>
    </div>
</div>
