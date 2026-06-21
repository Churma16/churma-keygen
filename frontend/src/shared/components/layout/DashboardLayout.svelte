<script>
    import { router } from 'svelte-spa-router';
    import { useQueryClient, useIsFetching } from '@tanstack/svelte-query';
    
    // Shared Layout
    import Sidebar from './Sidebar.svelte';
    import Header from './Header.svelte';
    import Toast from '../ui/Toast.svelte';
    import { toast } from '../../store/toastStore';
    import { licenseApi } from '../../../domains/license/api/licenseApi';

    // Svelte 5 children prop
    let { children } = $props();

    let isSidebarOpen = $state(false);

    // Derive activeTab from the router location property
    const activeTab = $derived(
        router.location === '/licenses' ? 'licenses' :
        router.location === '/clients' ? 'clients' :
        router.location === '/logs' ? 'logs' :
        router.location === '/settings' ? 'settings' :
        'overview'
    );

    // Breadcrumb mapping
    const breadcrumbPath = $derived(
        activeTab === 'overview' ? 'Ringkasan' :
        activeTab === 'licenses' ? 'Daftar Lisensi' :
        activeTab === 'clients' ? 'Klien Toko' :
        activeTab === 'settings' ? 'Pengaturan Akun' : 'Log Aktivasi'
    );

    // Global queries fetching count
    const isFetchingCount = useIsFetching();
    const isTabLoading = $derived(isFetchingCount.current > 0);

    const queryClient = useQueryClient();

    function handleRefresh() {
        queryClient.invalidateQueries();
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
</script>

<div class="min-h-screen bg-base-100 flex font-sans text-gray-700 overflow-x-hidden">
    <!-- LEFT SIDEBAR COMPONENT -->
    <Sidebar activeTab={activeTab} bind:isOpen={isSidebarOpen} />

    {#if isSidebarOpen}
        <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
        <div 
            class="fixed inset-0 z-30 bg-black/40 lg:hidden transition-opacity duration-300"
            on:click={() => isSidebarOpen = false}
        ></div>
    {/if}

    <!-- MAIN CANVAS -->
    <main class="flex-1 bg-base-100 min-h-screen flex flex-col min-w-0">
        <!-- HEADER COMPONENT -->
        <Header 
            breadcrumbPath={breadcrumbPath} 
            isLoadingData={isTabLoading} 
            on:copyPublicKey={copyPublicKey}
            on:refresh={handleRefresh}
            on:toggleSidebar={() => isSidebarOpen = !isSidebarOpen}
        />

        <!-- Page Content Area -->
        <div class="p-4 sm:p-8 flex-1 max-w-6xl w-full mx-auto" style="background-color: #faf8f5">
            {@render children()}
        </div>
    </main>
</div>

<!-- TOAST COMPONENT -->
<Toast message={$toast.message} type={$toast.type} />
