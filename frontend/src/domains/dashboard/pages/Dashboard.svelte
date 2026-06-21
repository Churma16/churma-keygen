<script>
    import { onMount } from 'svelte';
    import {
        Search,
        Plus,
        UserPlus
    } from 'lucide-svelte';
    
    import { authStore } from '../../auth/store/authStore';
    import { licenseStore } from '../../license/store/licenseStore';
    import { licenseApi } from '../../license/api/licenseApi';
    import { formatDate } from '../../../shared/utils/format';
    
    // Shared Layout
    import Sidebar from '../../../shared/components/layout/Sidebar.svelte';
    import Header from '../../../shared/components/layout/Header.svelte';
    import Toast from '../../../shared/components/ui/Toast.svelte';
    
    // Domain Components
    import StatsGrid from '../../license/components/StatsGrid.svelte';
    import QuickActions from '../../license/components/QuickActions.svelte';
    import LicenseTable from '../../license/components/LicenseTable.svelte';
    import ModalGenerator from '../../license/components/ModalGenerator.svelte';
    
    import ClientTable from '../../client/components/ClientTable.svelte';
    import CreateClientModal from '../../client/components/CreateClientModal.svelte';
    import EditClientModal from '../../client/components/EditClientModal.svelte';
    
    import LogTable from '../../log/components/LogTable.svelte';

    // State Management
    let activeTab = 'overview'; // 'overview' | 'licenses' | 'clients' | 'logs'
    let isSidebarOpen = false;

    // Search and Filter
    let searchQuery = '';
    let statusFilter = 'ALL';

    // Loading States
    let isLoadingData = false;
    let isSubmittingClient = false;

    // Toast notifications
    let toastMsg = '';
    let toastType = 'success'; // 'success' | 'error' | 'info'

    // Modals state
    let showGenModal = false;
    let selectedClientForLicense = null;
    let showCreateClientModal = false;

    // Edit client state
    let showEditClientModal = false;
    let editingClient = null;

    onMount(() => {
        if (!$authStore.isLoggedIn) {
            authStore.logout();
            return;
        }
        fetchDashboardData();
    });

    function showToast(msg, type = 'success') {
        toastMsg = msg;
        toastType = type;
        setTimeout(() => {
            toastMsg = '';
        }, 3500);
    }

    async function fetchDashboardData() {
        isLoadingData = true;
        try {
            await licenseStore.fetchAll();
        } catch (err) {
            showToast('Gagal memuat beberapa data dari server', 'error');
        } finally {
            isLoadingData = false;
        }
    }

    async function createClient(event) {
        const { name, ownerName, phone } = event.detail;
        isSubmittingClient = true;
        try {
            await licenseStore.createClient(name, ownerName, phone);
            showToast('Klien baru berhasil ditambahkan.');
            showCreateClientModal = false;
        } catch (e) {
            showToast(e.message || 'Gagal menambahkan klien.', 'error');
        } finally {
            isSubmittingClient = false;
        }
    }

    async function updateClient(event) {
        const { id, name, ownerName, phone } = event.detail;
        try {
            await licenseStore.updateClient(id, name, ownerName, phone);
            showToast('Identitas klien berhasil diperbarui.');
            showEditClientModal = false;
            editingClient = null;
        } catch (e) {
            showToast(e.message || 'Gagal memperbarui klien.', 'error');
        }
    }

    async function deleteClient(event) {
        const id = event.detail;
        if (!confirm('Apakah Anda yakin ingin menghapus klien ini secara permanen? Semua lisensi miliknya juga akan terhapus.')) return;
        try {
            await licenseStore.deleteClient(id);
            showToast('Klien berhasil dihapus (Soft Delete).');
        } catch (e) {
            showToast(e.message || 'Gagal menghapus klien.', 'error');
        }
    }

    async function updateLicenseStatus(event) {
        const { id, status } = event.detail;
        try {
            await licenseStore.updateLicenseStatus(id, status);
            showToast(`Status lisensi berhasil diubah menjadi ${status}.`);
        } catch (e) {
            showToast(e.message || 'Gagal mengubah status lisensi.', 'error');
        }
    }

    async function deleteLicense(event) {
        const id = event.detail;
        if (!confirm('Apakah Anda yakin ingin menghapus kunci lisensi ini?')) return;
        try {
            await licenseStore.deleteLicense(id);
            showToast('Lisensi berhasil dihapus (Soft Delete).');
        } catch (e) {
            showToast(e.message || 'Gagal menghapus lisensi.', 'error');
        }
    }

    async function copyPublicKey() {
        try {
            const text = await licenseApi.getPublicKey();
            await navigator.clipboard.writeText(text);
            showToast('Public Key RSA disalin! Siap dipasang di biner POS klien.', 'info');
        } catch (e) {
            showToast('Gagal mengunduh Public Key.', 'error');
        }
    }

    function openGenModal(event) {
        selectedClientForLicense = event.detail || null;
        showGenModal = true;
    }

    function openGenModalDirect(client = null) {
        selectedClientForLicense = client;
        showGenModal = true;
    }

    function openEditClientModal(event) {
        editingClient = event.detail;
        showEditClientModal = true;
    }

    // Breadcrumb mapping
    $: breadcrumbPath = activeTab === 'overview' ? 'Ringkasan' :
        activeTab === 'licenses' ? 'Daftar Lisensi' :
        activeTab === 'clients' ? 'Klien Toko' : 'Log Aktivasi';

    // Reactive filters
    $: filteredLicenses = ($licenseStore.licenses || []).filter(lic => {
        const clientName = lic.client_name || (lic.client ? lic.client.name : '');
        const matchesSearch =
            lic.license_code.toLowerCase().includes(searchQuery.toLowerCase()) ||
            clientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (lic.hardware_id && lic.hardware_id.toLowerCase().includes(searchQuery.toLowerCase()));
        const matchesStatus = statusFilter === 'ALL' || lic.status === statusFilter;
        return matchesSearch && matchesStatus;
    });

    $: filteredClients = ($licenseStore.clients || []).filter(c => {
        return c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (c.owner_name && c.owner_name.toLowerCase().includes(searchQuery.toLowerCase())) ||
            (c.phone && c.phone.includes(searchQuery));
    });

    $: filteredLogs = ($licenseStore.logs || []).filter(log => {
        return log.attempted_code.toLowerCase().includes(searchQuery.toLowerCase()) ||
            log.hardware_id_attempt.toLowerCase().includes(searchQuery.toLowerCase()) ||
            log.ip_address.includes(searchQuery) ||
            log.status.toLowerCase().includes(searchQuery.toLowerCase());
    });
</script>

<div class="min-h-screen bg-base-100 flex font-sans text-gray-700 overflow-x-hidden">
    <!-- LEFT SIDEBAR COMPONENT -->
    <Sidebar bind:activeTab bind:searchQuery bind:isOpen={isSidebarOpen} />

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
            isLoadingData={isLoadingData} 
            on:copyPublicKey={copyPublicKey}
            on:refresh={fetchDashboardData}
            on:toggleSidebar={() => isSidebarOpen = !isSidebarOpen}
        />

        <!-- Page Content Area -->
        <div class="p-4 sm:p-8 flex-1 max-w-6xl w-full mx-auto" style="background-color: #faf8f5">
            <div class="mb-8">
                {#if activeTab === 'overview'}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Selamat Datang, Admin</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Berikut adalah ringkasan performa lisensi dan audit sistem hari ini.</p>
                {:else if activeTab === 'licenses'}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Master Kunci Lisensi</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Buat, kelola status, batalkan, atau hapus kunci lisensi aplikasi POS klien.</p>
                {:else if activeTab === 'clients'}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Daftar Klien Toko</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Kelola identitas pemilik dan toko yang terintegrasi dengan lisensi keygen.</p>
                {:else}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Log Aktivasi & Audit Keamanan</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Pantau lalu lintas aktivasi sistem untuk mendeteksi cloning mesin atau pemalsuan HWID.</p>
                {/if}
            </div>

            <!-- STATS GRID COMPONENT (Visible on Overview page) -->
            {#if activeTab === 'overview'}
                <StatsGrid />

                <!-- Overview Dashboard Widgets -->
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
                    <!-- QUICK ACTIONS COMPONENT -->
                    <QuickActions 
                        on:addClient={() => showCreateClientModal = true}
                        on:printLicense={() => openGenModalDirect()}
                        on:copyPublicKey={copyPublicKey}
                    />

                    <!-- Right: Recent Activity Log summary -->
                    <div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm lg:col-span-2">
                        <div class="flex items-center justify-between mb-4">
                            <div>
                                <h3 class="font-bold text-primary text-sm mb-1">Aktivitas Aktivasi Terkini</h3>
                                <p class="text-[11px] text-gray-400 font-medium">Log upaya pencocokan serial kunci dan HWID oleh komputer kasir lokal.</p>
                            </div>
                            <button on:click={() => activeTab = 'logs'}
                                    class="btn btn-xs btn-ghost text-xs text-primary font-bold hover:bg-base-100">
                                Lihat Semua
                            </button>
                        </div>

                        <div class="overflow-x-auto">
                            <table class="table table-compact w-full text-xs">
                                <thead>
                                <tr class="text-gray-400 border-b border-base-300">
                                    <th class="py-2.5 font-bold">Klien</th>
                                    <th class="py-2.5 font-bold">IP</th>
                                    <th class="py-2.5 font-bold">Status</th>
                                    <th class="py-2.5 font-bold text-right">Waktu</th>
                                </tr>
                                </thead>
                                <tbody class="divide-y divide-base-300/65">
                                {#if ($licenseStore.logs || []).length === 0}
                                    <tr>
                                        <td colspan="4" class="text-center py-8 text-gray-400 font-medium">Tidak ada rekaman log.</td>
                                    </tr>
                                {:else}
                                    {#each $licenseStore.logs.slice(0, 5) as l}
                                        <tr>
                                            <td class="py-2.5 font-semibold text-primary">
                                                {#if l.license && l.license.client}
                                                    {l.license.client.name}
                                                {:else}
                                                    <span class="text-gray-400 font-normal italic">Kunci salah</span>
                                                {/if}
                                            </td>
                                            <td class="py-2.5 font-mono">{l.ip_address}</td>
                                            <td class="py-2.5">
                                                <span class="badge badge-sm font-bold text-[8px] px-1.5 py-0.5 border-none
                                                    {l.status === 'SUCCESS' ? 'badge-success text-primary bg-success/40' : ''}
                                                    {l.status === 'INVALID_KEY' ? 'badge-error text-red-700 bg-error/70' : ''}
                                                    {l.status === 'HWID_MISMATCH' ? 'badge-warning text-yellow-800 bg-yellow-100' : ''}
                                                    {l.status === 'SUSPENDED_KEY' ? 'badge-neutral text-amber-700 bg-amber-100' : ''}
                                                ">
                                                    {l.status}
                                                </span>
                                            </td>
                                            <td class="py-2.5 text-right text-gray-400 font-medium">{formatDate(l.created_at)}</td>
                                        </tr>
                                    {/each}
                                {/if}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            {/if}

            <!-- Data Tables (For tabs: licenses, clients, logs) -->
            {#if activeTab !== 'overview'}
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
                            {#if activeTab === 'licenses'}
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
                            {/if}

                            {#if activeTab === 'clients'}
                                <button on:click={() => showCreateClientModal = true}
                                        class="btn btn-sm btn-primary gap-1.5 rounded-md text-white text-xs font-bold">
                                    <UserPlus size={14}/>
                                    <span>Tambah Klien</span>
                                </button>
                            {:else if activeTab === 'licenses'}
                                <button on:click={() => openGenModalDirect()}
                                        class="btn btn-sm btn-primary gap-1.5 rounded-md text-white text-xs font-bold">
                                    <Plus size={14}/>
                                    <span>Cetak Lisensi</span>
                                </button>
                            {/if}
                        </div>
                    </div>

                    {#if $licenseStore.isLoading}
                        <div class="flex flex-col items-center justify-center h-80 gap-3">
                            <span class="loading loading-spinner loading-md text-primary"></span>
                            <span class="text-xs text-gray-400 font-semibold">Sinkronisasi data ke PostgreSQL...</span>
                        </div>
                    {:else}
                        {#if activeTab === 'licenses'}
                            <LicenseTable 
                                filteredLicenses={filteredLicenses}
                                on:updateStatus={updateLicenseStatus}
                                on:deleteLicense={deleteLicense}
                            />
                        {:else if activeTab === 'clients'}
                            <ClientTable 
                                filteredClients={filteredClients}
                                on:openGenModal={openGenModal}
                                on:editClient={openEditClientModal}
                                on:deleteClient={deleteClient}
                            />
                        {:else if activeTab === 'logs'}
                            <LogTable 
                                filteredLogs={filteredLogs}
                            />
                        {/if}
                    {/if}
                </div>
            {/if}
        </div>
    </main>
</div>

<!-- Modal: Create Client -->
{#if showCreateClientModal}
    <CreateClientModal 
        isSubmitting={isSubmittingClient}
        on:close={() => showCreateClientModal = false}
        on:submit={createClient}
    />
{/if}

<!-- Modal: Edit Client -->
{#if showEditClientModal && editingClient}
    <EditClientModal 
        client={editingClient}
        on:close={() => { showEditClientModal = false; editingClient = null; }}
        on:submit={updateClient}
    />
{/if}

<!-- License Generator Modal integration -->
{#if showGenModal}
    <ModalGenerator
            selectedClient={selectedClientForLicense}
            on:close={() => showGenModal = false}
    />
{/if}

<!-- TOAST COMPONENT -->
<Toast message={toastMsg} type={toastType} />
