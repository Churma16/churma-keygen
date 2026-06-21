<script>
  import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
  import Login from './domains/auth/pages/Login.svelte';
  import Dashboard from './domains/dashboard/pages/Dashboard.svelte';
  import { authStore } from './domains/auth/store/authStore';

  const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
  {#if $authStore.isLoggedIn}
    <Dashboard />
  {:else}
    <Login />
  {/if}
</QueryClientProvider>
