<script>
    import { createEventDispatcher } from 'svelte';
    import { ChevronRight, Copy, RefreshCw, Menu } from 'lucide-svelte';

    const dispatch = createEventDispatcher();

    export let breadcrumbPath = '';
    export let isLoadingData = false;
</script>

<header class="h-16 bg-white border-b border-base-300 px-4 sm:px-8 flex items-center justify-between sticky top-0 z-30 shadow-sm">
    <div class="flex items-center gap-2 text-xs font-semibold text-gray-400">
        <button on:click={() => dispatch('toggleSidebar')}
                class="btn btn-xs btn-ghost btn-circle w-8 h-8 lg:hidden border border-base-300 hover:bg-base-100 flex items-center justify-center mr-1">
            <Menu size={14} class="text-gray-700" />
        </button>
        <span class="hidden sm:inline">Dashboard</span>
        <ChevronRight size={12} class="hidden sm:inline"/>
        <span class="text-gray-700">{breadcrumbPath}</span>
    </div>

    <div class="flex items-center gap-2 sm:gap-3">
        <button on:click={() => dispatch('copyPublicKey')}
                class="btn btn-xs btn-outline btn-primary gap-1.5 h-8 px-2 sm:px-3 rounded-md text-xs font-bold transition duration-150">
            <Copy size={12}/>
            <span class="hidden md:inline">Salin RSA Public Key</span>
            <span class="md:hidden">RSA Key</span>
        </button>

        <button on:click={() => dispatch('refresh')} disabled={isLoadingData}
                class="btn btn-xs btn-ghost btn-circle w-8 h-8 border border-base-300 hover:bg-base-100 flex items-center justify-center shrink-0">
            <RefreshCw size={12} class={isLoadingData ? 'animate-spin' : ''}/>
        </button>
    </div>
</header>
