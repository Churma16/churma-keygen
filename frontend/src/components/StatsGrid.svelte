<script>
    import { Users, CheckCircle2, ShieldAlert, FileText } from 'lucide-svelte';
    import { licenseStore } from '../stores/licenseStore';
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
    <!-- Total Clients -->
    <div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
        <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-success/80 text-primary rounded-md flex items-center justify-center border border-success shrink-0">
                <Users size={20}/>
            </div>
            <div>
                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Total Klien Toko</div>
                <div class="text-2xl font-black mt-1 text-primary">
                    {$licenseStore.isLoading ? '...' : $licenseStore.stats.total_clients}
                </div>
            </div>
        </div>
    </div>

    <!-- Active Licenses -->
    <div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
        <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-success/80 text-primary rounded-md flex items-center justify-center border border-success shrink-0">
                <CheckCircle2 size={20}/>
            </div>
            <div>
                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Lisensi Aktif</div>
                <div class="text-2xl font-black mt-1 text-primary">
                    {$licenseStore.isLoading ? '...' : $licenseStore.stats.active_licenses}
                </div>
            </div>
        </div>
    </div>

    <!-- Suspended Licenses -->
    <div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
        <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-error text-secondary rounded-md flex items-center justify-center border border-red-200 shrink-0">
                <ShieldAlert size={20}/>
            </div>
            <div>
                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Ditangguhkan</div>
                <div class="text-2xl font-black mt-1 text-primary">
                    {$licenseStore.isLoading ? '...' : $licenseStore.stats.suspended_licenses}
                </div>
            </div>
        </div>
    </div>

    <!-- Unassigned Keys -->
    <div class="bg-white border border-base-300 rounded-lg p-6 shadow-sm flex items-center justify-between relative overflow-hidden">
        <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-md flex items-center justify-center border border-blue-100 shrink-0">
                <FileText size={20}/>
            </div>
            <div>
                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Belum Dipakai</div>
                <div class="text-2xl font-black mt-1 text-primary">
                    {$licenseStore.isLoading ? '...' : $licenseStore.stats.unassigned_licenses}
                </div>
            </div>
        </div>
    </div>
</div>
