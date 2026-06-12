<script>
    import { onMount } from 'svelte';
    import {
        Activity,
        CheckCircle2,
        Edit,
        Plus,
        Search,
        ShieldAlert,
        ShieldCheck,
        ShieldX,
        Trash2,
        UserPlus,
        X
    } from 'lucide-svelte';
    
    import { authStore } from '../stores/authStore';
    import { licenseStore } from '../stores/licenseStore';
    import { licenseApi } from '../lib/api/license';
    import { formatDate, formatQuota } from '../lib/format';
    import ModalGenerator from '../components/ModalGenerator.svelte';
    
    // Abstracted Components
    import Sidebar from '../components/Sidebar.svelte';
    import Header from '../components/Header.svelte';
    import StatsGrid from '../components/StatsGrid.svelte';
    import QuickActions from '../components/QuickActions.svelte';
    import Toast from '../components/Toast.svelte';

    // State Management
    let activeTab = 'overview'; // 'overview' | 'licenses' | 'clients' | 'logs'

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
    let newClientName = '';
    let newClientOwner = '';
    let newClientPhone = '';

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

    async function createClient() {
        if (!newClientName) return;
        isSubmittingClient = true;
        try {
            await licenseStore.createClient(newClientName, newClientOwner, newClientPhone);
            showToast('Klien baru berhasil ditambahkan.');
            newClientName = '';
            newClientOwner = '';
            newClientPhone = '';
            showCreateClientModal = false;
        } catch (e) {
            showToast(e.message || 'Gagal menambahkan klien.', 'error');
        } finally {
            isSubmittingClient = false;
        }
    }

    async function updateClient() {
        if (!editingClient || !editingClient.name) return;
        try {
            await licenseStore.updateClient(editingClient.id, editingClient.name, editingClient.owner_name, editingClient.phone);
            showToast('Identitas klien berhasil diperbarui.');
            showEditClientModal = false;
            editingClient = null;
        } catch (e) {
            showToast(e.message || 'Gagal memperbarui klien.', 'error');
        }
    }

    async function deleteClient(id) {
        if (!confirm('Apakah Anda yakin ingin menghapus klien ini secara permanen? Semua lisensi miliknya juga akan terhapus.')) return;
        try {
            await licenseStore.deleteClient(id);
            showToast('Klien berhasil dihapus (Soft Delete).');
        } catch (e) {
            showToast(e.message || 'Gagal menghapus klien.', 'error');
        }
    }

    async function updateLicenseStatus(id, newStatus) {
        try {
            await licenseStore.updateLicenseStatus(id, newStatus);
            showToast(`Status lisensi berhasil diubah menjadi ${newStatus}.`);
        } catch (e) {
            showToast(e.message || 'Gagal mengubah status lisensi.', 'error');
        }
    }

    async function deleteLicense(id) {
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

    function openGenModal(client = null) {
        selectedClientForLicense = client;
        showGenModal = true;
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

<div class="min-h-screen bg-base-100 flex font-sans text-gray-700">
    <!-- LEFT SIDEBAR COMPONENT -->
    <Sidebar bind:activeTab bind:searchQuery />

    <!-- MAIN CANVAS -->
    <main class="flex-1 bg-base-100 min-h-screen flex flex-col">
        <!-- HEADER COMPONENT -->
        <Header 
            breadcrumbPath={breadcrumbPath} 
            isLoadingData={isLoadingData} 
            on:copyPublicKey={copyPublicKey}
            on:refresh={fetchDashboardData}
        />

        <!-- Page Content Area -->
        <div class="p-8 flex-1 max-w-6xl w-full mx-auto" style="background-color: #faf8f5">
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
                        on:printLicense={() => openGenModal()}
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
                                <button on:click={() => openGenModal()}
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
                        <!-- SUB-PANEL: LICENSES -->
                        {#if activeTab === 'licenses'}
                            <div class="overflow-x-auto w-full">
                                <table class="table table-zebra w-full text-left">
                                    <thead>
                                    <tr class="bg-base-100/50 text-primary border-b border-base-300">
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">Klien Toko</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Kode Lisensi</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Hardware ID</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Batas Transaksi</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Status</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Masa Aktif</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">Aksi</th>
                                    </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-300">
                                    {#if filteredLicenses.length === 0}
                                        <tr>
                                            <td colspan="7" class="text-center py-16 text-gray-400 font-medium">Tidak ada lisensi terdaftar yang cocok.</td>
                                        </tr>
                                    {:else}
                                        {#each filteredLicenses as lic}
                                            <tr class="hover:bg-base-100/40 transition">
                                                <td class="py-4 font-semibold text-primary pl-6 max-w-[200px] truncate">
                                                    {lic.client_name || (lic.client ? lic.client.name : 'Unknown Client')}
                                                </td>
                                                <td class="py-4 font-mono font-bold text-gray-800 text-sm select-all">
                                                    {lic.license_code}
                                                </td>
                                                <td class="py-4 font-mono text-xs max-w-[150px] truncate text-gray-500">
                                                    {#if lic.hardware_id}
                                                        <span class="tooltip cursor-help" data-tip={lic.hardware_id}>
                                                            {lic.hardware_id.slice(0, 16)}...
                                                        </span>
                                                    {:else}
                                                        <span class="text-gray-400 italic">Belum terkunci</span>
                                                    {/if}
                                                </td>
                                                <td class="py-4 text-xs font-semibold">
                                                    {#if lic.trial_limit === 0 || lic.trial_limit === -1}
                                                        <span class="badge badge-success text-[10px] font-bold py-1.5 px-2 border-none text-primary bg-success/60">Unlimited</span>
                                                    {:else}
                                                        <span>{lic.trial_limit} Transaksi</span>
                                                    {/if}
                                                </td>
                                                <td class="py-4 text-xs">
                                                    {#if lic.status === 'ACTIVE'}
                                                        <span class="badge badge-success text-[9px] font-extrabold px-2 py-1 border-none text-primary bg-success/50">AKTIF</span>
                                                    {:else}
                                                        <span class="badge text-[9px] font-extrabold px-2 py-1 border-none
                                                            {lic.status === 'SUSPENDED' ? 'text-amber-800 bg-amber-100' : ''}
                                                            {lic.status === 'REVOKED' ? 'text-red-800 bg-error/70' : ''}
                                                            {lic.status === 'UNASSIGNED' ? 'text-gray-500 bg-gray-100' : ''}
                                                        ">
                                                            {lic.status}
                                                        </span>
                                                    {/if}
                                                </td>
                                                <td class="py-4 text-xs text-gray-400 font-semibold">
                                                    {#if lic.expires_at}
                                                        {formatDate(lic.expires_at)}
                                                    {:else}
                                                        <span class="text-gray-400 font-normal">Permanen</span>
                                                    {/if}
                                                </td>
                                                <td class="py-4 text-right pr-6">
                                                    <div class="flex justify-end gap-2">
                                                        {#if lic.status === 'ACTIVE'}
                                                            <button
                                                                    on:click={() => updateLicenseStatus(lic.id, 'SUSPENDED')}
                                                                    class="btn btn-xs btn-outline btn-warning gap-1 tooltip"
                                                                    data-tip="Tangguhkan Lisensi"
                                                            >
                                                                <ShieldAlert size={12}/>
                                                                <span>Tangguhkan</span>
                                                            </button>
                                                        {:else if lic.status === 'SUSPENDED'}
                                                            <button
                                                                    on:click={() => updateLicenseStatus(lic.id, 'ACTIVE')}
                                                                    class="btn btn-xs btn-outline btn-primary gap-1 tooltip"
                                                                    data-tip="Aktifkan Ulang"
                                                            >
                                                                <ShieldCheck size={12}/>
                                                                <span>Aktifkan</span>
                                                            </button>
                                                        {/if}

                                                        {#if lic.status !== 'REVOKED' && lic.status !== 'UNASSIGNED'}
                                                            <button
                                                                    on:click={() => updateLicenseStatus(lic.id, 'REVOKED')}
                                                                    class="btn btn-xs btn-outline btn-error gap-1 tooltip"
                                                                    data-tip="Batalkan Lisensi"
                                                            >
                                                                <ShieldX size={12}/>
                                                                <span>Batalkan</span>
                                                            </button>
                                                        {/if}

                                                        <button
                                                                on:click={() => deleteLicense(lic.id)}
                                                                class="w-6 h-6 border border-red-100 hover:bg-red-50 text-red-500 rounded-md flex items-center justify-center tooltip"
                                                                data-tip="Hapus Lisensi"
                                                        >
                                                            <Trash2 size={12}/>
                                                        </button>
                                                    </div>
                                                </td>
                                            </tr>
                                        {/each}
                                    {/if}
                                    </tbody>
                                </table>
                            </div>

                            <!-- SUB-PANEL: CLIENTS -->
                        {:else if activeTab === 'clients'}
                            <div class="overflow-x-auto w-full">
                                <table class="table table-zebra w-full text-left">
                                    <thead>
                                    <tr class="bg-base-100/50 text-primary border-b border-base-300">
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">Nama Klien / Perusahaan</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Nama Pemilik</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Nomor Telepon</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Jumlah Lisensi</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Terdaftar Pada</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">Aksi</th>
                                    </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-300">
                                    {#if filteredClients.length === 0}
                                        <tr>
                                            <td colspan="6" class="text-center py-16 text-gray-400 font-medium">Tidak ada klien toko terdaftar.</td>
                                        </tr>
                                    {:else}
                                        {#each filteredClients as c}
                                            <tr class="hover:bg-base-100/40 transition text-sm">
                                                <td class="py-4 font-semibold text-primary pl-6 select-all">{c.name}</td>
                                                <td class="py-4 font-medium text-gray-700">
                                                    {#if c.owner_name}
                                                        {c.owner_name}
                                                    {:else}
                                                        <span class="text-gray-400 italic">Kosong</span>
                                                    {/if}
                                                </td>
                                                <td class="py-4 select-all text-xs font-semibold text-gray-600">
                                                    {#if c.phone}
                                                        {c.phone}
                                                    {:else}
                                                        <span class="text-gray-400 italic">Kosong</span>
                                                    {/if}
                                                </td>
                                                <td class="py-4 text-xs">
                                                    <span class="badge badge-success text-[10px] font-bold py-1.5 px-2 border-none text-primary bg-success/60">
                                                        {c.licenses ? c.licenses.length : 0} kunci
                                                    </span>
                                                </td>
                                                <td class="py-4 text-xs text-gray-400 font-semibold">{formatDate(c.created_at)}</td>
                                                <td class="py-4 text-right pr-6">
                                                    <div class="flex justify-end gap-2">
                                                        <button
                                                                on:click={() => openGenModal(c)}
                                                                class="btn btn-xs btn-primary gap-1 rounded-md text-white"
                                                        >
                                                            <Plus size={12}/>
                                                            <span>Cetak Kunci</span>
                                                        </button>

                                                        <button
                                                                on:click={() => { editingClient = { ...c }; showEditClientModal = true; }}
                                                                class="w-6 h-6 border border-base-300 hover:bg-base-100 text-gray-600 rounded-md flex items-center justify-center tooltip"
                                                                data-tip="Edit Profil"
                                                        >
                                                            <Edit size={12}/>
                                                        </button>

                                                        <button
                                                                on:click={() => deleteClient(c.id)}
                                                                class="w-6 h-6 border border-red-100 hover:bg-red-50 text-red-500 rounded-md flex items-center justify-center tooltip"
                                                                data-tip="Hapus Klien"
                                                        >
                                                            <Trash2 size={12}/>
                                                        </button>
                                                    </div>
                                                </td>
                                            </tr>
                                        {/each}
                                    {/if}
                                    </tbody>
                                </table>
                            </div>

                            <!-- SUB-PANEL: LOGS -->
                        {:else if activeTab === 'logs'}
                            <div class="overflow-x-auto w-full">
                                <table class="table table-zebra w-full text-left">
                                    <thead>
                                    <tr class="bg-base-100/50 text-primary border-b border-base-300">
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">IP Address</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Kunci Lisensi</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Klien Terkait</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Hardware ID Percobaan</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Status Respon</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">Waktu Upaya</th>
                                    </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-300">
                                    {#if filteredLogs.length === 0}
                                        <tr>
                                            <td colspan="6" class="text-center py-16 text-gray-400 font-medium">Tidak ada log audit yang terekam.</td>
                                        </tr>
                                    {:else}
                                        {#each filteredLogs as l}
                                            <tr class="hover:bg-base-100/40 transition text-xs">
                                                <td class="py-3 pl-6 font-mono font-bold text-gray-700">{l.ip_address}</td>
                                                <td class="py-3 font-mono font-bold text-primary">{l.attempted_code}</td>
                                                <td class="py-3 font-semibold text-gray-800">
                                                    {#if l.license && l.license.client}
                                                        {l.license.client.name}
                                                    {:else}
                                                        <span class="text-gray-400 font-normal italic">Tidak terikat</span>
                                                    {/if}
                                                </td>
                                                <td class="py-3 font-mono text-[10px] max-w-[200px] truncate text-gray-500" title={l.hardware_id_attempt}>
                                                    {l.hardware_id_attempt}
                                                </td>
                                                <td class="py-3">
                                                    <span class="badge badge-sm font-bold text-[9px] px-1.5 py-0.5 border-none
                                                        {l.status === 'SUCCESS' ? 'text-primary bg-success/50' : ''}
                                                        {l.status === 'INVALID_KEY' ? 'text-red-800 bg-error/70' : ''}
                                                        {l.status === 'HWID_MISMATCH' ? 'text-yellow-800 bg-yellow-100' : ''}
                                                        {l.status === 'SUSPENDED_KEY' ? 'text-amber-800 bg-amber-100' : ''}
                                                    ">
                                                        {l.status}
                                                    </span>
                                                </td>
                                                <td class="py-3 text-right pr-6 text-gray-400 font-semibold">{formatDate(l.created_at)}</td>
                                            </tr>
                                        {/each}
                                    {/if}
                                    </tbody>
                                </table>
                            </div>
                        {/if}
                    {/if}
                </div>
            {/if}
        </div>
    </main>
</div>

<!-- Modal: Create Client -->
{#if showCreateClientModal}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="modal modal-open" on:click|self={() => showCreateClientModal = false}>
        <div class="modal-box bg-white border border-base-300 rounded-lg shadow-2xl relative p-6 max-w-md">
            <button on:click={() => showCreateClientModal = false}
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
                            bind:value={newClientName}
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
                            bind:value={newClientOwner}
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
                            bind:value={newClientPhone}
                            class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                    />
                </div>
            </div>

            <div class="modal-action gap-2 mt-6">
                <button on:click={() => showCreateClientModal = false}
                        class="btn btn-outline btn-sm rounded-md text-xs font-semibold h-9 px-4">Batal
                </button>
                <button
                        on:click={createClient}
                        disabled={!newClientName || isSubmittingClient}
                        class="btn btn-primary btn-sm text-white rounded-md text-xs font-bold h-9 px-4"
                >
                    {#if isSubmittingClient}
                        <span class="loading loading-spinner loading-xs mr-1"></span>
                    {/if}
                    Simpan Klien
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Modal: Edit Client -->
{#if showEditClientModal && editingClient}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="modal modal-open" on:click|self={() => { showEditClientModal = false; editingClient = null; }}>
        <div class="modal-box bg-white border border-base-300 rounded-lg shadow-2xl relative p-6 max-w-md">
            <button on:click={() => { showEditClientModal = false; editingClient = null; }}
                    class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4 text-gray-400 hover:text-gray-600">
                <X size={16}/>
            </button>

            <div class="flex items-center gap-3 mb-5">
                <div class="p-2.5 bg-success/80 text-primary rounded-md border border-success">
                    <Edit size={18}/>
                </div>
                <h3 class="font-bold text-lg text-primary">Edit Profil Klien</h3>
            </div>

            <div class="space-y-4">
                <div class="form-control">
                    <label class="label px-0.5 py-1" for="edit-name">
                        <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Toko</span>
                    </label>
                    <input
                            id="edit-name"
                            type="text"
                            bind:value={editingClient.name}
                            class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                    />
                </div>

                <div class="form-control">
                    <label class="label px-0.5 py-1" for="edit-owner">
                        <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Pemilik</span>
                    </label>
                    <input
                            id="edit-owner"
                            type="text"
                            bind:value={editingClient.owner_name}
                            class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                    />
                </div>

                <div class="form-control">
                    <label class="label px-0.5 py-1" for="edit-phone">
                        <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nomor Telepon</span>
                    </label>
                    <input
                            id="edit-phone"
                            type="text"
                            bind:value={editingClient.phone}
                            class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                    />
                </div>
            </div>

            <div class="modal-action gap-2 mt-6">
                <button on:click={() => { showEditClientModal = false; editingClient = null; }}
                        class="btn btn-outline btn-sm rounded-md text-xs font-semibold h-9 px-4">Batal
                </button>
                <button
                        on:click={updateClient}
                        disabled={!editingClient.name}
                        class="btn btn-primary btn-sm text-white rounded-md text-xs font-bold h-9 px-4"
                >
                    Simpan Perubahan
                </button>
            </div>
        </div>
    </div>
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
