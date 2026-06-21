<script>
    import { createEventDispatcher } from 'svelte';
    import { ShieldAlert, ShieldCheck, ShieldX, HelpCircle } from 'lucide-svelte';

    const dispatch = createEventDispatcher();

    export let title = 'Konfirmasi Tindakan';
    export let message = 'Apakah Anda yakin ingin melanjutkan tindakan ini?';
    export let confirmText = 'Ya, Lanjutkan';
    export let cancelText = 'Batal';
    export let variant = 'primary'; // 'primary' | 'warning' | 'error' | 'success'

    function handleConfirm() {
        dispatch('confirm');
    }

    function handleCancel() {
        dispatch('close');
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<div class="modal modal-open z-50 backdrop-blur-[2px] transition-all bg-black/40" on:click|self={handleCancel}>
    <div class="modal-box bg-white border border-base-300 rounded-lg shadow-2xl relative p-6 max-w-sm">
        
        <div class="flex flex-col items-center text-center mt-2">
            <!-- Icon container -->
            <div class="p-3.5 rounded-full mb-4 {
                variant === 'warning' ? 'bg-amber-50 text-amber-600 border border-amber-100' :
                variant === 'error' ? 'bg-rose-50 text-rose-600 border border-rose-100' :
                variant === 'success' ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' :
                variant === 'primary' ? 'bg-blue-50 text-blue-600 border border-blue-100' :
                'bg-gray-50 text-gray-600 border border-gray-100'
            }">
                {#if variant === 'warning'}
                    <ShieldAlert size={28} class="animate-pulse" />
                {:else if variant === 'error'}
                    <ShieldX size={28} />
                {:else if variant === 'success'}
                    <ShieldCheck size={28} />
                {:else}
                    <HelpCircle size={28} />
                {/if}
            </div>

            <h3 class="font-bold text-lg text-primary mb-2">{title}</h3>
            <p class="text-xs text-gray-500 font-semibold leading-relaxed px-2">
                {message}
            </p>
        </div>

        <div class="flex gap-2.5 mt-6 w-full justify-center">
            <button on:click={handleCancel}
                    class="btn btn-outline border-gray-300 hover:bg-gray-50 hover:text-gray-700 btn-sm rounded-md text-xs font-semibold h-9 px-5 flex-1">
                {cancelText}
            </button>
            <button
                    on:click={handleConfirm}
                    class="btn btn-sm text-white rounded-md text-xs font-bold h-9 px-5 flex-1 {
                        variant === 'warning' ? 'btn-warning bg-amber-500 hover:bg-amber-600 border-none' :
                        variant === 'error' ? 'btn-error bg-red-500 hover:bg-red-600 border-none' :
                        variant === 'success' ? 'btn-success bg-emerald-600 hover:bg-emerald-700 border-none' :
                        'btn-primary border-none'
                    }"
            >
                {confirmText}
            </button>
        </div>
    </div>
</div>
