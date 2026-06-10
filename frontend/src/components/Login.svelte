<script>
  import { createEventDispatcher } from 'svelte';
  import { KeyRound, User, Lock, AlertCircle, Loader2 } from 'lucide-svelte';

  const dispatch = createEventDispatcher();

  let username = '';
  let password = '';
  let error = '';
  let isLoading = false;

  async function handleLogin() {
    error = '';
    isLoading = true;

    try {
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      const body = await response.json();

      if (!response.ok) {
        throw new Error(body.meta?.message || body.error || 'Autentikasi gagal. Periksa kembali username & password.');
      }

      const data = body.data;

      // Save credentials & notify parent
      localStorage.setItem('admin_token', data.token);
      localStorage.setItem('admin_username', data.username);
      localStorage.setItem('admin_role', data.role);
      
      dispatch('loginSuccess', data);
    } catch (err) {
      error = err.message;
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-base-100 px-4 font-sans">
  <!-- Clean White Card -->
  <div class="w-full max-w-md bg-white shadow-xl rounded-2xl border border-base-300 p-8 relative overflow-hidden">
    <!-- Subtle color accent glows -->
    <div class="absolute -top-10 -right-10 w-40 h-40 bg-primary/5 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-10 -left-10 w-40 h-40 bg-secondary/5 rounded-full blur-3xl pointer-events-none"></div>

    <div class="flex flex-col items-center mb-8 relative z-10">
      <!-- Mint circle with Forest Green icon -->
      <div class="w-16 h-16 bg-success/60 text-primary rounded-2xl flex items-center justify-center mb-4 border border-success">
        <KeyRound size={28} />
      </div>
      <h1 class="text-2xl font-bold tracking-tight text-primary">Churma Keygen</h1>
      <p class="text-xs text-gray-500 mt-1 font-medium">Sistem Aktivasi Lisensi Administrator</p>
    </div>

    {#if error}
      <!-- Soft Red alert with dark red text -->
      <div class="alert alert-error mb-6 py-3 px-4 rounded-xl flex items-start text-xs border border-red-200/50 shadow-sm text-red-800">
        <AlertCircle size={16} class="mt-0.5 shrink-0 text-red-600" />
        <span class="font-semibold">{error}</span>
      </div>
    {/if}

    <form on:submit|preventDefault={handleLogin} class="space-y-5 relative z-10">
      <div class="form-control w-full">
        <label class="label px-0.5 py-1" for="username">
          <span class="label-text font-bold text-gray-700 text-xs uppercase tracking-wider">Username</span>
        </label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
            <User size={16} />
          </div>
          <input
            id="username"
            type="text"
            placeholder="Ketik username Anda..."
            bind:value={username}
            required
            disabled={isLoading}
            class="input input-bordered w-full pl-10 bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary transition duration-150 text-sm"
          />
        </div>
      </div>

      <div class="form-control w-full">
        <label class="label px-0.5 py-1" for="password">
          <span class="label-text font-bold text-gray-700 text-xs uppercase tracking-wider">Password</span>
        </label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
            <Lock size={16} />
          </div>
          <input
            id="password"
            type="password"
            placeholder="••••••••"
            bind:value={password}
            required
            disabled={isLoading}
            class="input input-bordered w-full pl-10 bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary transition duration-150 text-sm"
          />
        </div>
      </div>

      <button
        type="submit"
        disabled={isLoading}
        class="btn btn-primary w-full mt-4 h-12 rounded-xl text-white font-bold tracking-wide shadow-md shadow-primary/10 hover:shadow-primary/25 active:scale-[0.99] transition duration-150"
      >
        {#if isLoading}
          <Loader2 size={16} class="animate-spin mr-2" />
          Mencocokkan Kredensial...
        {:else}
          Masuk ke Dashboard
        {/if}
      </button>
    </form>
  </div>
</div>
