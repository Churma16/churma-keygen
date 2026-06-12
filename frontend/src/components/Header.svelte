<script>
    import { createEventDispatcher } from 'svelte';
    import { ChevronRight, Copy, RefreshCw } from 'lucide-svelte';

    const dispatch = createEventDispatcher();

    export let breadcrumbPath = '';
    export let isLoadingData = false;
</script>

<header class="h-16 bg-white border-b border-base-300 px-8 flex items-center justify-between sticky top-0 z-30 shadow-sm">
    <div class="flex items-center gap-2 text-xs font-semibold text-gray-400">
        <span>Dashboard</span>
        <ChevronRight size={12}/>
        <span class="text-gray-700">{breadcrumbPath}</span>
    </div>

    <div class="flex items-center gap-3">
        <button on:click={() => dispatch('copyPublicKey')}
                class="btn btn-xs btn-outline btn-primary gap-1.5 h-8 px-3 rounded-md text-xs font-bold transition duration-150">
            <Copy size={12}/>
            <span>Salin RSA Public Key</span>
        </button>

        <button on:click={() => dispatch('refresh')} disabled={isLoadingData}
                class="btn btn-xs btn-ghost btn-circle w-8 h-8 border border-base-300 hover:bg-base-100">
            <RefreshCw size={12} class={isLoadingData ? 'animate-spin' : ''}/>
        </button>
    </div>
</header>
