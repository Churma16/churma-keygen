<script>
  import { authStore } from '../stores/authStore';
  import { authApi } from '../lib/api/auth';
  import Input from '../components/ui/Input.svelte';
  import Button from '../components/ui/Button.svelte';
  import { KeyRound, User, Lock, AlertCircle, Loader2 } from 'lucide-svelte';

  let username = '';
  let password = '';
  let error = '';
  let isLoading = false;

  async function handleLogin() {
    error = '';
    isLoading = true;

    try {
      const res = await authApi.login(username, password);
      const data = res.data;
      
      // Save credentials & update state
      authStore.login(data.token, data.username, data.role);
    } catch (err) {
      error = err.message;
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-base-100 px-4 font-sans">
  <!-- Clean White Card -->
  <div class="w-full max-w-md bg-white shadow-xl rounded-lg border border-base-300 p-8 relative overflow-hidden">
    <!-- Subtle color accent glows -->
    <div class="absolute -top-10 -right-10 w-40 h-40 bg-primary/5 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-10 -left-10 w-40 h-40 bg-secondary/5 rounded-full blur-3xl pointer-events-none"></div>

    <div class="flex flex-col items-center mb-8 relative z-10">
      <!-- Mint circle with Forest Green icon -->
      <div class="w-16 h-16 bg-success/60 text-primary rounded-lg flex items-center justify-center mb-4 border border-success">
        <KeyRound size={28} />
      </div>
      <h1 class="text-2xl font-bold tracking-tight text-primary">Churma Keygen</h1>
      <p class="text-xs text-gray-500 mt-1 font-medium">Sistem Aktivasi Lisensi Administrator</p>
    </div>

    {#if error}
      <!-- Soft Red alert with dark red text -->
      <div class="alert alert-error mb-6 py-3 px-4 rounded-md flex items-start text-xs border border-red-200/50 shadow-sm text-red-800">
        <AlertCircle size={16} class="mt-0.5 shrink-0 text-red-600" />
        <span class="font-semibold">{error}</span>
      </div>
    {/if}

    <form on:submit|preventDefault={handleLogin} class="space-y-5 relative z-10">
      <Input
        id="username"
        type="text"
        placeholder="Ketik username Anda..."
        bind:value={username}
        required
        disabled={isLoading}
        label="Username"
      >
        <User size={16} slot="icon" />
      </Input>

      <Input
        id="password"
        type="password"
        placeholder="••••••••"
        bind:value={password}
        required
        disabled={isLoading}
        label="Password"
      >
        <Lock size={16} slot="icon" />
      </Input>

      <Button
        type="submit"
        disabled={isLoading}
        className="w-full mt-4"
      >
        {#if isLoading}
          <Loader2 size={16} class="animate-spin mr-2" />
          Mencocokkan Kredensial...
        {:else}
          Masuk ke Dashboard
        {/if}
      </Button>
    </form>
  </div>
</div>
