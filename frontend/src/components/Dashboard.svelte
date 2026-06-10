<script>
    import {onMount} from 'svelte';
    import {
        Activity,
        CheckCircle2,
        ChevronRight,
        Copy,
        Edit,
        FileText,
        Grid,
        Info,
        KeyRound,
        LogOut,
        Plus,
        RefreshCw,
        Search,
        Shield,
        ShieldAlert,
        ShieldCheck,
        ShieldX,
        Trash2,
        UserPlus,
        Users
    } from 'lucide-svelte';
    import LicenseGeneratorModal from './LicenseGeneratorModal.svelte';

    // State Management
    let activeTab = 'overview'; // 'overview' | 'licenses' | 'clients' | 'logs'
    let stats = {
        total_clients: 0,
        active_licenses: 0,
        suspended_licenses: 0,
        unassigned_licenses: 0,
        revoked_licenses: 0
    };
    let clients = [];
    let licenses = [];
    let logs = [];
    let token = localStorage.getItem('admin_token') || '';
    let username = localStorage.getItem('admin_username') || 'Admin';

    // Search and Filter
    let searchQuery = '';
    let statusFilter = 'ALL';

    // Loading States
    let isLoadingStats = false;
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
        if (!token) {
            window.location.reload();
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
            await Promise.all([
                fetchStats(),
                fetchLicenses(),
                fetchClients(),
                fetchLogs()
            ]);
        } catch (err) {
            showToast('Gagal memuat beberapa data dari server', 'error');
        } finally {
            isLoadingData = false;
        }
    }

    async function fetchStats() {
        isLoadingStats = true;
        try {
            const res = await fetch('/api/v1/admin/stats', {
                headers: {'Authorization': `Bearer ${token}`}
            });
            if (res.status === 401) handleSessionExpired();
            if (res.ok) {
                stats = await res.json();
            }
        } catch (e) {
            console.error(e);
        } finally {
            isLoadingStats = false;
        }
    }

    async function fetchLicenses() {
        try {
            const res = await fetch('/api/v1/admin/licenses', {
                headers: {'Authorization': `Bearer ${token}`}
            });
            if (res.ok) {
                licenses = await res.json();
            }
        } catch (e) {
            console.error(e);
        }
    }

    async function fetchClients() {
        try {
            const res = await fetch('/api/v1/admin/clients', {
                headers: {'Authorization': `Bearer ${token}`}
            });
            if (res.ok) {
                clients = await res.json();
            }
        } catch (e) {
            console.error(e);
        }
    }

    async function fetchLogs() {
        try {
            const res = await fetch('/api/v1/admin/logs', {
                headers: {'Authorization': `Bearer ${token}`}
            });
            if (res.ok) {
                logs = await res.json();
            }
        } catch (e) {
            console.error(e);
        }
    }

    function handleSessionExpired() {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_username');
        localStorage.removeItem('admin_role');
        window.location.reload();
    }

    async function createClient() {
        if (!newClientName) return;
        isSubmittingClient = true;
        try {
            const res = await fetch('/api/v1/admin/clients', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: newClientName,
                    owner_name: newClientOwner,
                    phone: newClientPhone
                })
            });
            if (res.ok) {
                showToast('Klien baru berhasil ditambahkan.');
                newClientName = '';
                newClientOwner = '';
                newClientPhone = '';
                showCreateClientModal = false;
                fetchClients();
                fetchStats();
            } else {
                const d = await res.json();
                showToast(d.error || 'Gagal menambahkan klien.', 'error');
            }
        } catch (e) {
            showToast('Koneksi terputus.', 'error');
        } finally {
            isSubmittingClient = false;
        }
    }

    async function updateClient() {
        if (!editingClient || !editingClient.name) return;
        try {
            const res = await fetch(`/api/v1/admin/clients/${editingClient.id}`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: editingClient.name,
                    owner_name: editingClient.owner_name,
                    phone: editingClient.phone
                })
            });
            if (res.ok) {
                showToast('Identitas klien berhasil diperbarui.');
                showEditClientModal = false;
                editingClient = null;
                fetchClients();
                fetchLicenses();
            } else {
                showToast('Gagal memperbarui klien.', 'error');
            }
        } catch (e) {
            showToast('Koneksi terputus.', 'error');
        }
    }

    async function deleteClient(id) {
        if (!confirm('Apakah Anda yakin ingin menghapus klien ini secara permanen? Semua lisensi miliknya juga akan terhapus.')) return;
        try {
            const res = await fetch(`/api/v1/admin/clients/${id}`, {
                method: 'DELETE',
                headers: {'Authorization': `Bearer ${token}`}
            });
            if (res.ok) {
                showToast('Klien berhasil dihapus (Soft Delete).');
                fetchClients();
                fetchLicenses();
                fetchStats();
            } else {
                showToast('Gagal menghapus klien.', 'error');
            }
        } catch (e) {
            showToast('Koneksi terputus.', 'error');
        }
    }

    async function updateLicenseStatus(id, newStatus) {
        try {
            const res = await fetch(`/api/v1/admin/licenses/${id}/status`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({status: newStatus})
            });
            if (res.ok) {
                showToast(`Status lisensi berhasil diubah menjadi ${newStatus}.`);
                fetchLicenses();
                fetchStats();
            } else {
                showToast('Gagal mengubah status lisensi.', 'error');
            }
        } catch (e) {
            showToast('Koneksi terputus.', 'error');
        }
    }

    async function deleteLicense(id) {
        if (!confirm('Apakah Anda yakin ingin menghapus kunci lisensi ini?')) return;
        try {
            const res = await fetch(`/api/v1/admin/licenses/${id}`, {
                method: 'DELETE',
                headers: {'Authorization': `Bearer ${token}`}
            });
            if (res.ok) {
                showToast('Lisensi berhasil dihapus (Soft Delete).');
                fetchLicenses();
                fetchStats();
            } else {
                showToast('Gagal menghapus lisensi.', 'error');
            }
        } catch (e) {
            showToast('Koneksi terputus.', 'error');
        }
    }

    async function copyPublicKey() {
        try {
            const res = await fetch('/api/v1/public-key');
            if (res.ok) {
                const text = await res.text();
                await navigator.clipboard.writeText(text);
                showToast('Public Key RSA disalin! Siap dipasang di biner POS klien.', 'info');
            } else {
                showToast('Gagal mengunduh Public Key.', 'error');
            }
        } catch (e) {
            showToast('Gagal terhubung dengan server.', 'error');
        }
    }

    function handleLogout() {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_username');
        localStorage.removeItem('admin_role');
        window.location.reload();
    }

    function openGenModal(client = null) {
        selectedClientForLicense = client;
        showGenModal = true;
    }

    function handleLicenseGenerated() {
        showGenModal = false;
        fetchLicenses();
        fetchStats();
    }

    function formatTime(t) {
        if (!t) return '-';
        return new Date(t).toLocaleString('id-ID', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // Breadcrumb mapping
    $: breadcrumbPath = activeTab === 'overview' ? 'Ringkasan' :
        activeTab === 'licenses' ? 'Daftar Lisensi' :
            activeTab === 'clients' ? 'Klien Toko' : 'Log Aktivasi';

    // Reactive filters
    $: filteredLicenses = licenses.filter(lic => {
        const matchesSearch =
            lic.license_code.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (lic.client && lic.client.name.toLowerCase().includes(searchQuery.toLowerCase())) ||
            (lic.hardware_id && lic.hardware_id.toLowerCase().includes(searchQuery.toLowerCase()));
        const matchesStatus = statusFilter === 'ALL' || lic.status === statusFilter;
        return matchesSearch && matchesStatus;
    });

    $: filteredClients = clients.filter(c => {
        return c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (c.owner_name && c.owner_name.toLowerCase().includes(searchQuery.toLowerCase())) ||
            (c.phone && c.phone.includes(searchQuery));
    });

    $: filteredLogs = logs.filter(log => {
        return log.attempted_code.toLowerCase().includes(searchQuery.toLowerCase()) ||
            log.hardware_id_attempt.toLowerCase().includes(searchQuery.toLowerCase()) ||
            log.ip_address.includes(searchQuery) ||
            log.status.toLowerCase().includes(searchQuery.toLowerCase());
    });
</script>

<div class="min-h-screen bg-base-100 flex font-sans text-gray-700">
    <!-- 1. LEFT SIDEBAR -->
    <aside class="w-64 bg-white border-r border-base-300 min-h-screen flex flex-col justify-between shrink-0">
        <div>
            <!-- Brand Logo / Title -->
            <div class="px-6 py-5 border-b border-base-300 flex items-center gap-2.5">
                <!-- 4-square grid logo in Coral Red -->
                <div class="grid grid-cols-2 gap-1 w-6 h-6 text-secondary shrink-0">
                    <div class="bg-current rounded-sm"></div>
                    <div class="bg-current rounded-sm"></div>
                    <div class="bg-current rounded-sm"></div>
                    <div class="bg-current rounded-sm"></div>
                </div>
                <div class="flex items-baseline font-bold text-lg tracking-tight">
                    <span class="text-primary">Churma</span>
                    <span class="text-secondary">Keygen</span>
                    <span class="text-secondary text-xs font-semibold ml-0.5">.ai</span>
                </div>
            </div>

            <!-- User Admin Profile Section -->
            <div class="px-6 py-5 flex items-center gap-3 border-b border-base-300 bg-base-100/50">
                <div class="avatar placeholder">
                    <div class="bg-primary/10 text-primary border border-primary/20 rounded-full w-10 h-10 font-bold text-sm">
                        {username.slice(0, 2).toUpperCase()}
                    </div>
                </div>
                <div>
                    <div class="font-bold text-sm text-primary select-all">{username}</div>
                    <div class="text-[10px] uppercase font-bold text-gray-400 tracking-wider">superadmin</div>
                </div>
            </div>

            <!-- Sidebar Navigation Menus -->
            <nav class="p-4 space-y-6">
                <!-- Menu Group: OPERASIONAL -->
                <div>
                    <div class="px-3 mb-2 text-[10px] font-bold text-gray-400 uppercase tracking-widest">Operasional
                    </div>
                    <ul class="space-y-1">
                        <li>
                            <button
                                    on:click={() => { activeTab = 'overview'; searchQuery = ''; }}
                                    class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold transition duration-150
                  {activeTab === 'overview' ? 'bg-primary text-white shadow-md' : 'text-gray-500 hover:bg-base-100 hover:text-primary'}"
                            >
                                <Grid size={16}/>
                                <span>Ringkasan</span>
                            </button>
                        </li>
                        <li>
                            <button
                                    on:click={() => { activeTab = 'licenses'; searchQuery = ''; }}
                                    class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold transition duration-150
                  {activeTab === 'licenses' ? 'bg-primary text-white shadow-md' : 'text-gray-500 hover:bg-base-100 hover:text-primary'}"
                            >
                                <KeyRound size={16}/>
                                <span>Daftar Lisensi</span>
                            </button>
                        </li>
                    </ul>
                </div>

                <!-- Menu Group: LAPORAN & DATA -->
                <div>
                    <div class="px-3 mb-2 text-[10px] font-bold text-gray-400 uppercase tracking-widest">Laporan &
                        Data
                    </div>
                    <ul class="space-y-1">
                        <li>
                            <button
                                    on:click={() => { activeTab = 'clients'; searchQuery = ''; }}
                                    class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold transition duration-150
                  {activeTab === 'clients' ? 'bg-primary text-white shadow-md' : 'text-gray-500 hover:bg-base-100 hover:text-primary'}"
                            >
                                <Users size={16}/>
                                <span>Klien Toko</span>
                            </button>
                        </li>
                        <li>
                            <button
                                    on:click={() => { activeTab = 'logs'; searchQuery = ''; }}
                                    class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold transition duration-150
                  {activeTab === 'logs' ? 'bg-primary text-white shadow-md' : 'text-gray-500 hover:bg-base-100 hover:text-primary'}"
                            >
                                <Activity size={16}/>
                                <span>Log Aktivasi</span>
                            </button>
                        </li>
                    </ul>
                </div>
            </nav>
        </div>

        <!-- Sidebar Footer section -->
        <div class="p-4 border-t border-base-300 space-y-3 bg-base-100/30">
            <!-- "Masa Aktif Server" / Active Info Badge (Crown style) -->
            <div class="bg-primary text-white rounded-2xl p-4 shadow-inner relative overflow-hidden flex flex-col gap-1">
                <div class="absolute -right-6 -bottom-6 text-white/10 pointer-events-none">
                    <Shield size={64}/>
                </div>
                <div class="flex items-center gap-2">
                    <Shield size={16} class="text-secondary"/>
                    <span class="font-bold text-xs tracking-wide uppercase">Gateway Lisensi</span>
                </div>
                <div class="text-[10px] text-gray-200/90 font-medium leading-relaxed mt-1">
                    Server Kripto RSA-2048 Aktif & Mengamankan Aplikasi POS.
                </div>
            </div>

            <!-- Logout button -->
            <button
                    on:click={handleLogout}
                    class="w-full btn btn-sm btn-ghost hover:bg-red-50 text-red-600 hover:text-red-700 flex justify-start gap-2.5 px-4 h-10 rounded-xl border border-transparent hover:border-red-100 text-xs font-bold transition duration-150"
            >
                <LogOut size={14}/>
                <span>Keluar Sistem</span>
            </button>
        </div>
    </aside>

    <!-- 2. MAIN CANVAS -->
    <main class="flex-1 bg-base-100 min-h-screen flex flex-col">
        <!-- Top breadcrumbs and quick actions bar -->
        <header class="h-16 bg-white border-b border-base-300 px-8 flex items-center justify-between sticky top-0 z-30 shadow-sm">
            <!-- Breadcrumbs -->
            <div class="flex items-center gap-2 text-xs font-semibold text-gray-400">
                <span>Dashboard</span>
                <ChevronRight size={12}/>
                <span class="text-gray-700">{breadcrumbPath}</span>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-3">
                <!-- Copy RSA Public Key -->
                <button on:click={copyPublicKey}
                        class="btn btn-xs btn-outline btn-primary gap-1.5 h-8 px-3 rounded-lg text-xs font-bold transition duration-150">
                    <Copy size={12}/>
                    <span>Salin RSA Public Key</span>
                </button>

                <!-- Refresh button -->
                <button on:click={fetchDashboardData} disabled={isLoadingData}
                        class="btn btn-xs btn-ghost btn-circle w-8 h-8 border border-base-300 hover:bg-base-100">
                    <RefreshCw size={12} class={isLoadingData ? 'animate-spin' : ''}/>
                </button>
            </div>
        </header>

        <!-- Page Content Area -->
        <div class="p-8 flex-1 max-w-6xl w-full mx-auto" style="background-color: #faf8f5">
            <!-- Header banner based on active navigation tab -->
            <div class="mb-8">
                {#if activeTab === 'overview'}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Selamat Datang, Admin</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Berikut adalah ringkasan performa lisensi dan
                        audit sistem hari ini.</p>
                {:else if activeTab === 'licenses'}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Master Kunci Lisensi</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Buat, kelola status, batalkan, atau hapus kunci
                        lisensi aplikasi POS klien.</p>
                {:else if activeTab === 'clients'}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Daftar Klien Toko</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Kelola identitas pemilik dan toko yang
                        terintegrasi dengan lisensi keygen.</p>
                {:else}
                    <h2 class="text-2xl font-bold text-primary tracking-tight">Log Aktivasi & Audit Keamanan</h2>
                    <p class="text-xs text-gray-500 mt-1 font-medium">Pantau lalu lintas aktivasi sistem untuk
                        mendeteksi cloning mesin atau pemalsuan HWID.</p>
                {/if}
            </div>

            <!-- Stats Cards (Visible on Overview/Overview pages) -->
            {#if activeTab === 'overview'}
                <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                    <!-- Total Clients -->
                    <div class="bg-white border border-base-300 rounded-2xl p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
                        <div class="flex items-center gap-4">
                            <!-- Success Mint circle background -->
                            <div class="w-12 h-12 bg-success/80 text-primary rounded-xl flex items-center justify-center border border-success shrink-0">
                                <Users size={20}/>
                            </div>
                            <div>
                                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Total Klien
                                    Toko
                                </div>
                                <div class="text-2xl font-black mt-1 text-primary">
                                    {isLoadingStats ? '...' : stats.total_clients}
                                </div>
                            </div>
                        </div>
                        <!-- Green upward trend badge -->
                        <div class="badge badge-success text-[10px] font-bold py-1 px-1.5 border-none text-primary bg-success/40 shrink-0 self-start mt-0.5">
                            +12%
                        </div>
                    </div>

                    <!-- Active Licenses -->
                    <div class="bg-white border border-base-300 rounded-2xl p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
                        <div class="flex items-center gap-4">
                            <div class="w-12 h-12 bg-success/80 text-primary rounded-xl flex items-center justify-center border border-success shrink-0">
                                <CheckCircle2 size={20}/>
                            </div>
                            <div>
                                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Lisensi
                                    Aktif
                                </div>
                                <div class="text-2xl font-black mt-1 text-primary">
                                    {isLoadingStats ? '...' : stats.active_licenses}
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Suspended Licenses -->
                    <div class="bg-white border border-base-300 rounded-2xl p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
                        <div class="flex items-center gap-4">
                            <!-- Coral Red background for Warning -->
                            <div class="w-12 h-12 bg-error text-secondary rounded-xl flex items-center justify-center border border-red-200 shrink-0">
                                <ShieldAlert size={20}/>
                            </div>
                            <div>
                                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Ditangguhkan
                                </div>
                                <div class="text-2xl font-black mt-1 text-primary">
                                    {isLoadingStats ? '...' : stats.suspended_licenses}
                                </div>
                            </div>
                        </div>
                        <!-- Coral upward trend badge -->
                        <div class="badge badge-error text-[10px] font-bold py-1 px-1.5 border-none text-secondary bg-error/70 shrink-0 self-start mt-0.5">
                            +5%
                        </div>
                    </div>

                    <!-- Unassigned Keys -->
                    <div class="bg-white border border-base-300 rounded-2xl p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
                        <div class="flex items-center gap-4">
                            <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center border border-blue-100 shrink-0">
                                <FileText size={20}/>
                            </div>
                            <div>
                                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Belum
                                    Dipakai
                                </div>
                                <div class="text-2xl font-black mt-1 text-primary">
                                    {isLoadingStats ? '...' : stats.unassigned_licenses}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Overview Dashboard layout widgets: Quick statistics + Recents -->
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    <!-- Left: Quick Links / Actions -->
                    <div class="bg-white border border-base-300 rounded-2xl p-6 shadow-sm flex flex-col justify-between">
                        <div>
                            <h3 class="font-bold text-primary text-sm mb-1">Pintasan Aksi</h3>
                            <p class="text-[11px] text-gray-400 font-medium mb-4">Gunakan pintasan cepat ini untuk
                                mengelola lisensi klien dengan instan.</p>

                            <div class="space-y-2">
                                <button on:click={() => showCreateClientModal = true}
                                        class="btn btn-sm btn-outline btn-primary w-full gap-2 justify-start h-10 rounded-xl text-xs font-bold">
                                    <UserPlus size={14}/>
                                    <span>Tambah Klien Baru</span>
                                </button>
                                <button on:click={() => openGenModal()}
                                        class="btn btn-sm btn-primary w-full gap-2 justify-start h-10 rounded-xl text-xs font-bold text-white">
                                    <Plus size={14}/>
                                    <span>Cetak Kunci Lisensi</span>
                                </button>
                                <button on:click={copyPublicKey}
                                        class="btn btn-sm btn-outline btn-ghost w-full gap-2 justify-start h-10 rounded-xl text-xs font-semibold text-gray-600 border-base-300 hover:bg-base-100">
                                    <Copy size={14}/>
                                    <span>Salin RSA Public Key</span>
                                </button>
                            </div>
                        </div>

                        <!-- Server Status Widget -->
                        <div class="mt-6 border-t border-base-300 pt-4 flex items-center gap-3">
              <span class="relative flex h-2.5 w-2.5 shrink-0">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
                <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-success"></span>
              </span>
                            <div class="text-[11px] font-semibold text-gray-500">
                                Koneksi database PostgreSQL terhubung dan aman.
                            </div>
                        </div>
                    </div>

                    <!-- Right: Recent Activity Log summary -->
                    <div class="bg-white border border-base-300 rounded-2xl p-6 shadow-sm lg:col-span-2">
                        <div class="flex items-center justify-between mb-4">
                            <div>
                                <h3 class="font-bold text-primary text-sm mb-1">Aktivitas Aktivasi Terkini</h3>
                                <p class="text-[11px] text-gray-400 font-medium">Log upaya pencocokan serial kunci dan
                                    HWID oleh komputer kasir lokal.</p>
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
                                {#if logs.length === 0}
                                    <tr>
                                        <td colspan="4" class="text-center py-8 text-gray-400 font-medium">Tidak ada
                                            rekaman log.
                                        </td>
                                    </tr>
                                {:else}
                                    {#each logs.slice(0, 5) as l}
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
                          <span class="badge badge-sm font-bold text-[8px] px-1.5 py-0.5
                            {l.status === 'SUCCESS' ? 'badge-success text-primary bg-success/40 border-none' : ''}
                            {l.status === 'INVALID_KEY' ? 'badge-error text-red-700 bg-error/70 border-none' : ''}
                            {l.status === 'HWID_MISMATCH' ? 'badge-warning text-yellow-800 bg-yellow-100 border-none' : ''}
                            {l.status === 'SUSPENDED_KEY' ? 'badge-neutral text-amber-700 bg-amber-100 border-none' : ''}
                          ">
                            {l.status}
                          </span>
                                            </td>
                                            <td class="py-2.5 text-right text-gray-400 font-medium">{formatTime(l.created_at)}</td>
                                        </tr>
                                    {/each}
                                {/if}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            {/if}

            <!-- Clean tables and search elements (For tabs: licenses, clients, logs) -->
            {#if activeTab !== 'overview'}
                <div class="bg-white border border-base-300 rounded-2xl shadow-sm overflow-hidden min-h-[400px]">
                    <!-- Inner Controls header -->
                    <div class="px-6 py-4 border-b border-base-300 bg-base-100/30 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                        <div class="relative grow max-w-sm">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                                <Search size={14}/>
                            </div>
                            <input
                                    type="text"
                                    placeholder="Cari data di tabel ini..."
                                    bind:value={searchQuery}
                                    class="input input-sm input-bordered pl-9 bg-white border-base-300 text-gray-800 w-full rounded-lg focus:outline-none focus:border-primary text-xs"
                            />
                        </div>

                        <div class="flex items-center gap-3 self-end sm:self-auto">
                            <!-- Filter status only for licenses tab -->
                            {#if activeTab === 'licenses'}
                                <select
                                        bind:value={statusFilter}
                                        class="select select-sm select-bordered bg-white border-base-300 text-gray-700 rounded-lg focus:outline-none text-xs"
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
                                        class="btn btn-sm btn-primary gap-1.5 rounded-lg text-white text-xs font-bold">
                                    <UserPlus size={14}/>
                                    <span>Tambah Klien</span>
                                </button>
                            {:else if activeTab === 'licenses'}
                                <!-- Forest green / primary generator button -->
                                <button on:click={() => openGenModal()}
                                        class="btn btn-sm btn-primary gap-1.5 rounded-lg text-white text-xs font-bold">
                                    <Plus size={14}/>
                                    <span>Cetak Lisensi</span>
                                </button>
                            {/if}
                        </div>
                    </div>

                    <!-- Loading Data Spinner -->
                    {#if isLoadingData}
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
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">Klien
                                            Toko
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Kode Lisensi
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Hardware ID
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Batas
                                            Transaksi
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Status</th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Masa Aktif
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">
                                            Aksi
                                        </th>
                                    </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-300">
                                    {#if filteredLicenses.length === 0}
                                        <tr>
                                            <td colspan="7" class="text-center py-16 text-gray-400 font-medium">
                                                Tidak ada lisensi terdaftar yang cocok.
                                            </td>
                                        </tr>
                                    {:else}
                                        {#each filteredLicenses as lic}
                                            <tr class="hover:bg-base-100/40 transition">
                                                <td class="py-4 font-semibold text-primary pl-6 max-w-[200px] truncate">
                                                    {lic.client ? lic.client.name : 'Unknown Client'}
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
                                                    {#if lic.trial_limit === 0}
                                                        <!-- Mint green badge for Unlimited -->
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
                                                        {formatTime(lic.expires_at)}
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

                                                        <!-- Round icon delete button -->
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
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">Nama
                                            Klien / Perusahaan
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Nama Pemilik
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Nomor
                                            Telepon
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Jumlah
                                            Lisensi
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Terdaftar
                                            Pada
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">
                                            Aksi
                                        </th>
                                    </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-300">
                                    {#if filteredClients.length === 0}
                                        <tr>
                                            <td colspan="6" class="text-center py-16 text-gray-400 font-medium">
                                                Tidak ada klien toko terdaftar.
                                            </td>
                                        </tr>
                                    {:else}
                                        {#each filteredClients as c}
                                            <tr class="hover:bg-base-100/40 transition text-sm">
                                                <td class="py-4 font-semibold text-primary pl-6 select-all">
                                                    {c.name}
                                                </td>
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
                                                    <span class="badge badge-success text-[10px] font-bold py-1.5 px-2 border-none text-primary bg-success/60">{c.licenses ? c.licenses.length : 0}
                                                        kunci</span>
                                                </td>
                                                <td class="py-4 text-xs text-gray-400 font-semibold">
                                                    {formatTime(c.created_at)}
                                                </td>
                                                <td class="py-4 text-right pr-6">
                                                    <div class="flex justify-end gap-2">
                                                        <button
                                                                on:click={() => openGenModal(c)}
                                                                class="btn btn-xs btn-primary gap-1 rounded-lg text-white"
                                                        >
                                                            <Plus size={12}/>
                                                            <span>Cetak Kunci</span>
                                                        </button>

                                                        <!-- Edit Button -->
                                                        <button
                                                                on:click={() => { editingClient = { ...c }; showEditClientModal = true; }}
                                                                class="w-6 h-6 border border-base-300 hover:bg-base-100 text-gray-600 rounded-md flex items-center justify-center tooltip"
                                                                data-tip="Edit Profil"
                                                        >
                                                            <Edit size={12}/>
                                                        </button>

                                                        <!-- Delete Button -->
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
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider pl-6">IP
                                            Address
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Kunci
                                            Lisensi
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Klien
                                            Terkait
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Hardware ID
                                            Percobaan
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider">Status
                                            Respon
                                        </th>
                                        <th class="py-3.5 font-bold text-[10px] uppercase tracking-wider text-right pr-6">
                                            Waktu Upaya
                                        </th>
                                    </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-300">
                                    {#if filteredLogs.length === 0}
                                        <tr>
                                            <td colspan="6" class="text-center py-16 text-gray-400 font-medium">
                                                Tidak ada log audit yang terekam.
                                            </td>
                                        </tr>
                                    {:else}
                                        {#each filteredLogs as l}
                                            <tr class="hover:bg-base-100/40 transition text-xs">
                                                <td class="py-3 pl-6 font-mono font-bold text-gray-700">
                                                    {l.ip_address}
                                                </td>
                                                <td class="py-3 font-mono font-bold text-primary">
                                                    {l.attempted_code}
                                                </td>
                                                <td class="py-3 font-semibold text-gray-800">
                                                    {#if l.license && l.license.client}
                                                        {l.license.client.name}
                                                    {:else}
                                                        <span class="text-gray-400 font-normal italic">Tidak terikat</span>
                                                    {/if}
                                                </td>
                                                <td class="py-3 font-mono text-[10px] max-w-[200px] truncate text-gray-500"
                                                    title={l.hardware_id_attempt}>
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
                                                <td class="py-3 text-right pr-6 text-gray-400 font-semibold">
                                                    {formatTime(l.created_at)}
                                                </td>
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

<!-- Modal: Create Client (Light theme) -->
{#if showCreateClientModal}
    <div class="modal modal-open">
        <div class="modal-box bg-white border border-base-300 rounded-2xl shadow-2xl relative p-6 max-w-md">
            <!-- Close Button -->
            <button on:click={() => showCreateClientModal = false}
                    class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4 text-gray-400 hover:text-gray-600">
                <X size={16}/>
            </button>

            <div class="flex items-center gap-3 mb-5">
                <div class="p-2.5 bg-success/80 text-primary rounded-xl border border-success">
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
                            class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
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
                            class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
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
                            class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
                    />
                </div>
            </div>

            <div class="modal-action gap-2 mt-6">
                <button on:click={() => showCreateClientModal = false}
                        class="btn btn-outline btn-sm rounded-lg text-xs font-semibold h-9 px-4">Batal
                </button>
                <button
                        on:click={createClient}
                        disabled={!newClientName || isSubmittingClient}
                        class="btn btn-primary btn-sm text-white rounded-lg text-xs font-bold h-9 px-4"
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

<!-- Modal: Edit Client (Light theme) -->
{#if showEditClientModal && editingClient}
    <div class="modal modal-open">
        <div class="modal-box bg-white border border-base-300 rounded-2xl shadow-2xl relative p-6 max-w-md">
            <button on:click={() => { showEditClientModal = false; editingClient = null; }}
                    class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4 text-gray-400 hover:text-gray-600">
                <X size={16}/>
            </button>

            <div class="flex items-center gap-3 mb-5">
                <div class="p-2.5 bg-success/80 text-primary rounded-xl border border-success">
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
                            class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
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
                            class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
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
                            class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
                    />
                </div>
            </div>

            <div class="modal-action gap-2 mt-6">
                <button on:click={() => { showEditClientModal = false; editingClient = null; }}
                        class="btn btn-outline btn-sm rounded-lg text-xs font-semibold h-9 px-4">Batal
                </button>
                <button
                        on:click={updateClient}
                        disabled={!editingClient.name}
                        class="btn btn-primary btn-sm text-white rounded-lg text-xs font-bold h-9 px-4"
                >
                    Simpan Perubahan
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- License Generator Modal integration -->
{#if showGenModal}
    <LicenseGeneratorModal
            {clients}
            {token}
            selectedClient={selectedClientForLicense}
            on:close={() => showGenModal = false}
            on:success={handleLicenseGenerated}
    />
{/if}

<!-- Toast Component -->
{#if toastMsg}
    <div class="toast toast-end toast-bottom z-50">
        <div class="alert shadow-2xl py-3 px-5 rounded-xl text-xs border border-base-300 flex items-center gap-2
      {toastType === 'success' ? 'bg-success text-primary border-success' : ''}
      {toastType === 'error' ? 'bg-error text-red-800 border-red-200' : ''}
      {toastType === 'info' ? 'bg-blue-50 text-blue-800 border-blue-100' : ''}
    ">
            {#if toastType === 'success'}
                <CheckCircle2 size={14} class="text-primary"/>
            {:else}
                <Info size={14}/>
            {/if}
            <span class="font-bold">{toastMsg}</span>
        </div>
    </div>
{/if}
