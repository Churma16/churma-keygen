<script>
    import { Grid, KeyRound, Users, Activity, Shield, LogOut } from 'lucide-svelte';
    import { authStore } from '../stores/authStore';

    export let activeTab = 'overview';
    export let searchQuery = '';
    export let isOpen = false;

    function handleTabChange(tab) {
        activeTab = tab;
        searchQuery = '';
        isOpen = false;
    }
</script>

<aside class="fixed inset-y-0 left-0 z-40 w-64 bg-white border-r border-base-300 min-h-screen flex flex-col justify-between shrink-0 transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-auto {isOpen ? 'translate-x-0' : '-translate-x-full'}">
    <div>
        <!-- Brand Logo / Title -->
        <div class="px-6 py-5 border-b border-base-300 flex items-center gap-2.5">
            <div class="grid grid-cols-2 gap-1 w-6 h-6 text-secondary shrink-0">
                <div class="bg-current rounded-sm"></div>
                <div class="bg-current rounded-sm"></div>
                <div class="bg-current rounded-sm"></div>
                <div class="bg-current rounded-sm"></div>
            </div>
            <div class="flex items-baseline font-bold text-lg tracking-tight">
                <span class="text-primary">Churma</span>
                <span class="text-secondary">Keygen</span>
            </div>
        </div>

        <!-- User Admin Profile Section -->
        <div class="px-6 py-5 flex items-center gap-3 border-b border-base-300 bg-base-100/50">
            <div class="avatar placeholder">
                <div class="bg-primary/10 text-primary border border-primary/20 rounded-full w-10 h-10 font-bold text-sm">
                    {$authStore.username ? $authStore.username.slice(0, 2).toUpperCase() : 'AD'}
                </div>
            </div>
            <div>
                <div class="font-bold text-sm text-primary select-all">{$authStore.username || 'Admin'}</div>
                <div class="text-[10px] uppercase font-bold text-gray-400 tracking-wider">superadmin</div>
            </div>
        </div>

        <!-- Sidebar Navigation Menus -->
        <nav class="p-4 space-y-6">
            <!-- Menu Group: OPERASIONAL -->
            <div>
                <div class="px-3 mb-2 text-[10px] font-bold text-gray-400 uppercase tracking-widest">Operasional</div>
                <ul class="space-y-1">
                    <li>
                        <button
                                on:click={() => handleTabChange('overview')}
                                class="w-full flex items-center gap-3 px-4 py-2.5 rounded-md text-sm font-semibold transition duration-150
              {activeTab === 'overview' ? 'bg-primary text-white shadow-md' : 'text-gray-500 hover:bg-base-100 hover:text-primary'}"
                        >
                            <Grid size={16}/>
                            <span>Ringkasan</span>
                        </button>
                    </li>
                    <li>
                        <button
                                on:click={() => handleTabChange('licenses')}
                                class="w-full flex items-center gap-3 px-4 py-2.5 rounded-md text-sm font-semibold transition duration-150
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
                <div class="px-3 mb-2 text-[10px] font-bold text-gray-400 uppercase tracking-widest">Laporan & Data</div>
                <ul class="space-y-1">
                    <li>
                        <button
                                on:click={() => handleTabChange('clients')}
                                class="w-full flex items-center gap-3 px-4 py-2.5 rounded-md text-sm font-semibold transition duration-150
              {activeTab === 'clients' ? 'bg-primary text-white shadow-md' : 'text-gray-500 hover:bg-base-100 hover:text-primary'}"
                        >
                            <Users size={16}/>
                            <span>Klien Toko</span>
                        </button>
                    </li>
                    <li>
                        <button
                                on:click={() => handleTabChange('logs')}
                                class="w-full flex items-center gap-3 px-4 py-2.5 rounded-md text-sm font-semibold transition duration-150
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
        <div class="bg-primary text-white rounded-lg p-4 shadow-inner relative overflow-hidden flex flex-col gap-1">
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

        <button
                on:click={authStore.logout}
                class="w-full btn btn-sm btn-ghost hover:bg-red-50 text-red-600 hover:text-red-700 flex justify-start gap-2.5 px-4 h-10 rounded-md border border-transparent hover:border-red-100 text-xs font-bold transition duration-150"
        >
            <LogOut size={14}/>
            <span>Keluar Sistem</span>
        </button>
    </div>
</aside>
