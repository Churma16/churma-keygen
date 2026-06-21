<script>
    import { Search, Plus } from 'lucide-svelte';
    
    // TanStack Query Hooks
    import { 
        useLicensesQuery, 
        useUpdateLicenseStatusMutation, 
        useDeleteLicenseMutation 
    } from '../store/licenseQueries';
    
    // Shared Layout & UI
    import ConfirmationModal from '../../../shared/components/ui/ConfirmationModal.svelte';
    import { toast } from '../../../shared/store/toastStore';

    // Domain Components
    import LicenseTable from '../components/LicenseTable.svelte';
    import ModalGenerator from '../components/ModalGenerator.svelte';

    // State Management
    let searchQuery = $state('');
    let statusFilter = $state('ALL');

    // Modals state
    let showGenModal = $state(false);
    let selectedClientForLicense = $state(null);

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
    const licensesQuery = useLicensesQuery();
    const updateLicenseStatusMutation = useUpdateLicenseStatusMutation();
    const deleteLicenseMutation = useDeleteLicenseMutation();

    // Reactive filters — $derived ensures re-computation when TanStack Query data arrives
    const filteredLicenses = $derived(
        (licensesQuery.data || []).filter(lic => {
            const clientName = lic.client_name || (lic.client ? lic.client.name : '');
            const matchesSearch =
                lic.license_code.toLowerCase().includes(searchQuery.toLowerCase()) ||
                clientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
                (lic.hardware_id && lic.hardware_id.toLowerCase().includes(searchQuery.toLowerCase()));
            const matchesStatus = statusFilter === 'ALL' || lic.status === statusFilter;
            return matchesSearch && matchesStatus;
        })
    );

    function openGenModalDirect(client = null) {
        selectedClientForLicense = client;
        showGenModal = true;
    }

    async function handleUpdateLicenseStatus(event) {
        const { id, status } = event.detail;

        let title = 'Ubah Status Lisensi?';
        let message = 'Apakah Anda yakin ingin mengubah status lisensi ini?';
        let confirmText = 'Ya, Ubah';
        let variant = 'primary';

        if (status === 'SUSPENDED') {
            title = 'Tangguhkan Lisensi?';
            message = 'Apakah Anda yakin ingin menangguhkan lisensi ini? Aplikasi POS klien sementara tidak akan dapat memverifikasi lisensi ini.';
            confirmText = 'Ya, Tangguhkan';
            variant = 'warning';
        } else if (status === 'ACTIVE') {
            title = 'Aktifkan Ulang Lisensi?';
            message = 'Apakah Anda yakin ingin mengaktifkan kembali lisensi ini? Aplikasi POS klien akan berfungsi normal.';
            confirmText = 'Ya, Aktifkan';
            variant = 'success';
        } else if (status === 'REVOKED') {
            title = 'Batalkan Lisensi Permanen?';
            message = 'Apakah Anda yakin ingin membatalkan lisensi ini? Tindakan ini akan memutus akses POS klien secara permanen dan tidak dapat diaktifkan kembali.';
            confirmText = 'Ya, Batalkan';
            variant = 'error';
        }

        triggerConfirmation({
            title,
            message,
            confirmText,
            variant,
            onConfirm: async () => {
                try {
                    await updateLicenseStatusMutation.mutateAsync({ id, status });
                    toast.show(`Status lisensi berhasil diperbarui menjadi ${status}.`);
                } catch (e) {
                    toast.show(e.message || 'Gagal memperbarui status lisensi.', 'error');
                }
                closeConfirmation();
            }
        });
    }

    async function handleDeleteLicense(event) {
        const id = event.detail;
        triggerConfirmation({
            title: 'Hapus Kunci Lisensi?',
            message: 'Apakah Anda yakin ingin menghapus kunci lisensi ini secara permanen dari sistem?',
            confirmText: 'Ya, Hapus',
            variant: 'error',
            onConfirm: async () => {
                try {
                    await deleteLicenseMutation.mutateAsync(id);
                    toast.show('Lisensi berhasil dihapus (Soft Delete).');
                } catch (e) {
                    toast.show(e.message || 'Gagal menghapus lisensi.', 'error');
                }
                closeConfirmation();
            }
        });
    }
</script>

<div class="mb-8">
    <h2 class="text-2xl font-bold text-primary tracking-tight">Master Kunci Lisensi</h2>
    <p class="text-xs text-gray-500 mt-1 font-medium">Buat, kelola status, batalkan, atau hapus kunci lisensi aplikasi POS klien.</p>
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
            <select
                    bind:value={statusFilter}
                    class="select select-sm select-bordered bg-white border-base-300 text-gray-700 rounded-md focus:outline-none text-xs"
            >
                <option value="ALL">Semua Status</option>
                <option value="UNASSIGNED">Belum Dipakai</option>
                <option value="ACTIVE">Aktif</option>
                <option value="SUSPENDED">Ditangguhkan</option>
                <option value="REVOKED">Dibatalkan</option>
            </select>

            <button on:click={() => openGenModalDirect()}
                    class="btn btn-sm btn-primary gap-1.5 rounded-md text-white text-xs font-bold">
                <Plus size={14}/>
                <span>Cetak Lisensi</span>
            </button>
        </div>
    </div>

    {#if licensesQuery.isFetching && (licensesQuery.data || []).length === 0}
        <div class="flex flex-col items-center justify-center h-80 gap-3">
            <span class="loading loading-spinner loading-md text-primary"></span>
            <span class="text-xs text-gray-400 font-semibold">Sinkronisasi data ke PostgreSQL...</span>
        </div>
    {:else}
        <LicenseTable 
            filteredLicenses={filteredLicenses}
            on:updateStatus={handleUpdateLicenseStatus}
            on:deleteLicense={handleDeleteLicense}
        />
    {/if}
</div>

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
