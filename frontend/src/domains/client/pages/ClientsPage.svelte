<script>
    import { Search, UserPlus } from 'lucide-svelte';
    
    // TanStack Query Hooks
    import { 
        useClientsQuery, 
        useCreateClientMutation, 
        useUpdateClientMutation, 
        useDeleteClientMutation 
    } from '../store/clientQueries';
    
    // Shared Layout & UI
    import ConfirmationModal from '../../../shared/components/ui/ConfirmationModal.svelte';
    import { toast } from '../../../shared/store/toastStore';

    // Domain Components
    import ClientTable from '../components/ClientTable.svelte';
    import CreateClientModal from '../components/CreateClientModal.svelte';
    import EditClientModal from '../components/EditClientModal.svelte';
    import ModalGenerator from '../../license/components/ModalGenerator.svelte';

    // State Management
    let searchQuery = $state('');

    // Modals state
    let showGenModal = $state(false);
    let selectedClientForLicense = $state(null);
    let showCreateClientModal = $state(false);
    
    // Edit client state
    let showEditClientModal = $state(false);
    let editingClient = $state(null);

    // Confirmation modal state
    let showConfirmModal = $state(false);
    let confirmModalConfig = $state({
        title: '',
        message: '',
        confirmText: '',
        variant: 'primary',
        onConfirm: null
    });

    function triggerConfirmation(config) {
        confirmModalConfig = config;
        showConfirmModal = true;
    }

    function closeConfirmation() {
        showConfirmModal = false;
    }

    // Initialize TanStack Query Hooks
    const clientsQuery = useClientsQuery();
    const createClientMutation = useCreateClientMutation();
    const updateClientMutation = useUpdateClientMutation();
    const deleteClientMutation = useDeleteClientMutation();

    // Reactive filters — $derived ensures re-computation when TanStack Query data arrives
    const filteredClients = $derived(
        (clientsQuery.data || []).filter(cli => 
            cli.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            cli.owner_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            cli.phone.includes(searchQuery)
        )
    );

    function openGenModal(event) {
        selectedClientForLicense = event.detail || null;
        showGenModal = true;
    }

    function openEditClientModal(event) {
        editingClient = event.detail;
        showEditClientModal = true;
    }

    async function handleCreateClient(event) {
        const { name, ownerName, phone } = event.detail;
        try {
            await createClientMutation.mutateAsync({ name, ownerName, phone });
            toast.show('Klien baru berhasil ditambahkan.');
            showCreateClientModal = false;
        } catch (e) {
            toast.show(e.message || 'Gagal menambahkan klien.', 'error');
        }
    }

    async function handleUpdateClient(event) {
        const { id, name, ownerName, phone } = event.detail;
        try {
            await updateClientMutation.mutateAsync({ id, name, ownerName, phone });
            toast.show('Identitas klien berhasil diperbarui.');
            showEditClientModal = false;
            editingClient = null;
        } catch (e) {
            toast.show(e.message || 'Gagal memperbarui klien.', 'error');
        }
    }

    async function handleDeleteClient(event) {
        const id = event.detail;
        triggerConfirmation({
            title: 'Hapus Klien & Semua Lisensi?',
            message: 'Apakah Anda yakin ingin menghapus klien ini secara permanen? Semua lisensi milik klien ini juga akan ikut terhapus dari sistem.',
            confirmText: 'Ya, Hapus Klien',
            variant: 'error',
            onConfirm: async () => {
                try {
                    await deleteClientMutation.mutateAsync(id);
                    toast.show('Klien berhasil dihapus (Soft Delete).');
                } catch (e) {
                    toast.show(e.message || 'Gagal menghapus klien.', 'error');
                }
                closeConfirmation();
            }
        });
    }
</script>

<div class="mb-8">
    <h2 class="text-2xl font-bold text-primary tracking-tight">Daftar Klien Toko</h2>
    <p class="text-xs text-gray-500 mt-1 font-medium">Kelola identitas pemilik dan toko yang terintegrasi dengan lisensi keygen.</p>
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

        <div class="flex items-center gap-3 self-end sm:self-auto">
            <button on:click={() => showCreateClientModal = true}
                    class="btn btn-sm btn-primary gap-1.5 rounded-md text-white text-xs font-bold">
                <UserPlus size={14}/>
                <span>Tambah Klien</span>
            </button>
        </div>
    </div>

    {#if clientsQuery.isFetching && (clientsQuery.data || []).length === 0}
        <div class="flex flex-col items-center justify-center h-80 gap-3">
            <span class="loading loading-spinner loading-md text-primary"></span>
            <span class="text-xs text-gray-400 font-semibold">Sinkronisasi data ke PostgreSQL...</span>
        </div>
    {:else}
        <ClientTable 
            filteredClients={filteredClients}
            on:openGenModal={openGenModal}
            on:editClient={openEditClientModal}
            on:deleteClient={handleDeleteClient}
        />
    {/if}
</div>

<!-- Modal: Create Client -->
{#if showCreateClientModal}
    <CreateClientModal 
        isSubmitting={createClientMutation.isPending}
        on:close={() => showCreateClientModal = false}
        on:submit={handleCreateClient}
    />
{/if}

<!-- Modal: Edit Client -->
{#if showEditClientModal && editingClient}
    <EditClientModal 
        client={editingClient}
        isSubmitting={updateClientMutation.isPending}
        on:close={() => { showEditClientModal = false; editingClient = null; }}
        on:submit={handleUpdateClient}
    />
{/if}

<!-- License Generator Modal integration -->
{#if showGenModal}
    <ModalGenerator
            selectedClient={selectedClientForLicense}
            on:close={() => showGenModal = false}
    />
{/if}

<!-- Modal: Confirmation Modal -->
{#if showConfirmModal}
    <ConfirmationModal
        title={confirmModalConfig.title}
        message={confirmModalConfig.message}
        confirmText={confirmModalConfig.confirmText}
        variant={confirmModalConfig.variant}
        on:confirm={confirmModalConfig.onConfirm}
        on:close={closeConfirmation}
    />
{/if}
