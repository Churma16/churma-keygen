<script>
    import { createEventDispatcher } from 'svelte';
    import { User, Lock, Save, AlertCircle, ShieldCheck, Phone, Mail } from 'lucide-svelte';
    import { authStore } from '../../auth/store/authStore';
    import { useUpdateProfileMutation } from '../../auth/store/authQueries';
    import { useGetSetting, useUpdateSettingMutation } from '../store/settingQueries';

    const dispatch = createEventDispatcher();

    let username = $authStore.username || '';
    let currentPassword = '';
    let newPassword = '';
    let confirmPassword = '';

    let errorMsg = '';
    let successMsg = '';

    const updateProfileMutation = useUpdateProfileMutation();

    // Reset messages when input changes
    $: if (username || currentPassword || newPassword || confirmPassword) {
        errorMsg = '';
    }

    // Support Contact settings logic
    const getWhatsappQuery = useGetSetting('contact_whatsapp');
    const getEmailQuery = useGetSetting('contact_email');
    const updateWhatsappMutation = useUpdateSettingMutation();
    const updateEmailMutation = useUpdateSettingMutation();

    let whatsappPhone = '';
    let emailContact = '';
    let hasLoadedWhatsapp = false;
    let hasLoadedEmail = false;

    let contactErrorMsg = '';
    let contactSuccessMsg = '';

    $: if (getWhatsappQuery.data && !hasLoadedWhatsapp) {
        whatsappPhone = getWhatsappQuery.data.value || '';
        hasLoadedWhatsapp = true;
    }

    $: if (getEmailQuery.data && !hasLoadedEmail) {
        emailContact = getEmailQuery.data.value || '';
        hasLoadedEmail = true;
    }

    $: if (whatsappPhone || emailContact) {
        contactErrorMsg = '';
    }

    $: isContactSaving = updateWhatsappMutation.isPending || updateEmailMutation.isPending;
    $: isContactLoading = getWhatsappQuery.isPending || getEmailQuery.isPending;

    async function handleUpdateContact() {
        contactErrorMsg = '';
        contactSuccessMsg = '';

        try {
            // Use two separate mutation instances to avoid stuck isPending state
            await Promise.all([
                updateWhatsappMutation.mutateAsync({
                    key: 'contact_whatsapp',
                    value: whatsappPhone
                }),
                updateEmailMutation.mutateAsync({
                    key: 'contact_email',
                    value: emailContact
                })
            ]);
            contactSuccessMsg = 'Kontak dukungan berhasil diperbarui!';
        } catch (err) {
            contactErrorMsg = err.response?.data?.message || err.message || 'Gagal memperbarui kontak dukungan.';
        }
    }


    async function handleUpdateProfile() {
        if (!username) {
            errorMsg = 'Username wajib diisi.';
            return;
        }
        if (!currentPassword) {
            errorMsg = 'Password saat ini wajib diisi untuk verifikasi keamanan.';
            return;
        }
        if (newPassword && newPassword !== confirmPassword) {
            errorMsg = 'Konfirmasi password baru tidak cocok.';
            return;
        }

        errorMsg = '';
        successMsg = '';

        try {
            await updateProfileMutation.mutateAsync({
                username,
                current_password: currentPassword,
                new_password: newPassword
            });

            successMsg = 'Profil berhasil diperbarui! Mengalihkan ke halaman login dalam 2 detik...';
            
            // Clear inputs
            currentPassword = '';
            newPassword = '';
            confirmPassword = '';

            // Force logout after 2 seconds
            setTimeout(() => {
                authStore.logout();
            }, 2000);
        } catch (err) {
            errorMsg = err.response?.data?.message || err.message || 'Gagal memperbarui profil.';
        }
    }
</script>

<div class="bg-white border border-base-300 rounded-lg shadow-sm max-w-2xl w-full mx-auto overflow-hidden">
    <!-- Header Card -->
    <div class="px-6 py-5 border-b border-base-300 bg-base-100/30 flex items-center gap-3">
        <div class="p-2.5 bg-primary/10 text-primary rounded-md border border-primary/20">
            <User size={18}/>
        </div>
        <div>
            <h3 class="font-bold text-base text-primary">Kredensial Administrator</h3>
            <p class="text-xs text-gray-400 font-medium">Ubah nama pengguna atau password untuk masuk ke sistem.</p>
        </div>
    </div>

    <!-- Form Body -->
    <div class="p-6 space-y-5">
        {#if errorMsg}
            <div class="alert alert-error rounded-md flex items-start gap-2.5 text-xs text-red-800 bg-red-50 border border-red-200 p-3.5">
                <AlertCircle size={16} class="shrink-0 mt-0.5" />
                <span>{errorMsg}</span>
            </div>
        {/if}

        {#if successMsg}
            <div class="alert alert-success rounded-md flex items-start gap-2.5 text-xs text-emerald-800 bg-emerald-50 border border-emerald-200 p-3.5">
                <ShieldCheck size={16} class="shrink-0 mt-0.5" />
                <span>{successMsg}</span>
            </div>
        {/if}

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Username Input -->
            <div class="form-control md:col-span-2">
                <label class="label px-0.5 py-1" for="settings-username">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Username</span>
                </label>
                <div class="relative">
                    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                        <User size={16}/>
                    </span>
                    <input
                            id="settings-username"
                            type="text"
                            placeholder="Username"
                            bind:value={username}
                            class="input input-bordered w-full pl-10 bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                            disabled={updateProfileMutation.isPending}
                    />
                </div>
            </div>

            <!-- Current Password (for confirmation) -->
            <div class="form-control md:col-span-2">
                <label class="label px-0.5 py-1" for="settings-current-password">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Password Saat Ini</span>
                </label>
                <div class="relative">
                    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                        <Lock size={16}/>
                    </span>
                    <input
                            id="settings-current-password"
                            type="password"
                            placeholder="Konfirmasi password saat ini"
                            bind:value={currentPassword}
                            class="input input-bordered w-full pl-10 bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                            disabled={updateProfileMutation.isPending}
                    />
                </div>
                <div class="px-0.5 py-1">
                    <span class="text-[10px] text-gray-400 font-medium">Wajib diisi untuk memverifikasi identitas Anda sebelum menyimpan perubahan.</span>
                </div>
            </div>

            <div class="border-t border-base-300 md:col-span-2 my-2"></div>

            <!-- New Password -->
            <div class="form-control">
                <label class="label px-0.5 py-1" for="settings-new-password">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Password Baru</span>
                </label>
                <div class="relative">
                    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                        <Lock size={16}/>
                    </span>
                    <input
                            id="settings-new-password"
                            type="password"
                            placeholder="Kosongkan jika tidak ingin diubah"
                            bind:value={newPassword}
                            class="input input-bordered w-full pl-10 bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                            disabled={updateProfileMutation.isPending}
                    />
                </div>
            </div>

            <!-- Confirm New Password -->
            <div class="form-control">
                <label class="label px-0.5 py-1" for="settings-confirm-password">
                    <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Konfirmasi Password Baru</span>
                </label>
                <div class="relative">
                    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                        <Lock size={16}/>
                    </span>
                    <input
                            id="settings-confirm-password"
                            type="password"
                            placeholder="Ulangi password baru"
                            bind:value={confirmPassword}
                            class="input input-bordered w-full pl-10 bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                            disabled={updateProfileMutation.isPending}
                    />
                </div>
            </div>
        </div>
    </div>

    <!-- Footer Card with Actions -->
    <div class="px-6 py-4 border-t border-base-300 bg-base-100/10 flex justify-end gap-3">
        <button
                on:click={handleUpdateProfile}
                disabled={updateProfileMutation.isPending || !username || !currentPassword}
                class="btn btn-primary btn-sm text-white rounded-md text-xs font-bold h-9 px-5 flex items-center gap-2"
        >
            {#if updateProfileMutation.isPending}
                <span class="loading loading-spinner loading-xs"></span>
            {:else}
                <Save size={14}/>
            {/if}
            Simpan Perubahan
        </button>
    </div>
</div>

<div class="bg-white border border-base-300 rounded-lg shadow-sm max-w-2xl w-full mx-auto overflow-hidden mt-6">
    <!-- Header Card -->
    <div class="px-6 py-5 border-b border-base-300 bg-base-100/30 flex items-center gap-3">
        <div class="p-2.5 bg-primary/10 text-primary rounded-md border border-primary/20">
            <Phone size={18}/>
        </div>
        <div>
            <h3 class="font-bold text-base text-primary">Kontak Dukungan (WhatsApp & Email)</h3>
            <p class="text-xs text-gray-400 font-medium">Ubah nomor WhatsApp dan email dukungan yang akan ditampilkan pada aplikasi klien.</p>
        </div>
    </div>

    <!-- Form Body -->
    <div class="p-6 space-y-5">
        {#if isContactLoading}
            <div class="flex justify-center py-4">
                <span class="loading loading-spinner loading-md text-primary"></span>
            </div>

        {:else}
            {#if contactErrorMsg}
                <div class="alert alert-error rounded-md flex items-start gap-2.5 text-xs text-red-800 bg-red-50 border border-red-200 p-3.5">
                    <AlertCircle size={16} class="shrink-0 mt-0.5" />
                    <span>{contactErrorMsg}</span>
                </div>
            {/if}

            {#if contactSuccessMsg}
                <div class="alert alert-success rounded-md flex items-start gap-2.5 text-xs text-emerald-800 bg-emerald-50 border border-emerald-200 p-3.5">
                    <ShieldCheck size={16} class="shrink-0 mt-0.5" />
                    <span>{contactSuccessMsg}</span>
                </div>
            {/if}

            <div class="grid grid-cols-1 gap-4">
                <!-- WhatsApp Input -->
                <div class="form-control">
                    <label class="label px-0.5 py-1" for="settings-whatsapp">
                        <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Nomor WhatsApp</span>
                    </label>
                    <div class="relative">
                        <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                            <Phone size={16}/>
                        </span>
                        <input
                                id="settings-whatsapp"
                                type="text"
                                placeholder="Contoh: 6281234567890"
                                bind:value={whatsappPhone}
                                class="input input-bordered w-full pl-10 bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                                disabled={isContactSaving}
                        />
                    </div>
                    <div class="px-0.5 py-1">
                        <span class="text-[10px] text-gray-400 font-medium">Gunakan kode negara tanpa tanda "+" (misal: 62812xxxx). Nomor ini akan muncul di aplikasi klien saat tombol hubungi diklik.</span>
                    </div>
                </div>

                <div class="border-t border-base-300 my-2"></div>

                <!-- Email Input -->
                <div class="form-control">
                    <label class="label px-0.5 py-1" for="settings-email">
                        <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Email Dukungan</span>
                    </label>
                    <div class="relative">
                        <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                            <Mail size={16}/>
                        </span>
                        <input
                                id="settings-email"
                                type="email"
                                placeholder="Contoh: support@churma.com"
                                bind:value={emailContact}
                                class="input input-bordered w-full pl-10 bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary text-sm"
                                disabled={isContactSaving}
                        />
                    </div>
                    <div class="px-0.5 py-1">
                        <span class="text-[10px] text-gray-400 font-medium">Email dukungan yang dapat dihubungi oleh klien jika terjadi masalah aktivasi.</span>
                    </div>
                </div>
            </div>
        {/if}
    </div>

    <!-- Footer Card with Actions -->
    <div class="px-6 py-4 border-t border-base-300 bg-base-100/10 flex justify-end gap-3">
        <button
                on:click={handleUpdateContact}
                disabled={isContactLoading || isContactSaving}
                class="btn btn-primary btn-sm text-white rounded-md text-xs font-bold h-9 px-5 flex items-center gap-2"
        >
            {#if isContactSaving}
                <span class="loading loading-spinner loading-xs"></span>
            {:else}
                <Save size={14}/>
            {/if}
            Simpan Kontak
        </button>
    </div>
</div>
