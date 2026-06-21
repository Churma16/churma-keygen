<script>
    import { Search } from 'lucide-svelte';
    
    // TanStack Query Hooks
    import { useLogsQuery } from '../store/logQueries';

    // Domain Components
    import LogTable from '../components/LogTable.svelte';

    // State Management
    let searchQuery = $state('');

    // Initialize TanStack Query Hooks
    const logsQuery = useLogsQuery();

    // Reactive filters — $derived ensures re-computation when TanStack Query data arrives
    const filteredLogs = $derived(
        (logsQuery.data || []).filter(l => {
            const clientName = l.license && l.license.client ? l.license.client.name : '';
            const matchesSearch =
                clientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
                l.ip_address.toLowerCase().includes(searchQuery.toLowerCase()) ||
                l.status.toLowerCase().includes(searchQuery.toLowerCase()) ||
                (l.license && l.license.license_code.toLowerCase().includes(searchQuery.toLowerCase()));
            return matchesSearch;
        })
    );
</script>

<div class="mb-8">
    <h2 class="text-2xl font-bold text-primary tracking-tight">Log Aktivasi & Audit Keamanan</h2>
    <p class="text-xs text-gray-500 mt-1 font-medium">Pantau lalu lintas aktivasi sistem untuk mendeteksi cloning mesin atau pemalsuan HWID.</p>
</div>

<div class="bg-white border border-base-300 rounded-lg shadow-sm overflow-hidden min-h-[400px]">
    <div class="px-6 py-4 border-b border-base-300 bg-base-100/30 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <div class="relative grow max-w-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                <Search size={14}/>
            </div>
            <input
                    type="text"
                    placeholder="Cari data di tabel ini..."
                    bind:value={searchQuery}
                    class="input input-sm input-bordered pl-9 bg-white border-base-300 text-gray-800 w-full rounded-md focus:outline-none focus:border-primary text-xs"
             />
        </div>
    </div>

    {#if logsQuery.isFetching && (logsQuery.data || []).length === 0}
        <div class="flex flex-col items-center justify-center h-80 gap-3">
            <span class="loading loading-spinner loading-md text-primary"></span>
            <span class="text-xs text-gray-400 font-semibold">Sinkronisasi data ke PostgreSQL...</span>
        </div>
    {:else}
        <LogTable 
            filteredLogs={filteredLogs}
        />
    {/if}
</div>
