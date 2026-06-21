<script>
  import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
  import Router from 'svelte-spa-router';
  
  import { authStore } from './domains/auth/store/authStore';
  import Login from './domains/auth/pages/Login.svelte';

  // Layout & Pages
  import DashboardLayout from './shared/components/layout/DashboardLayout.svelte';
  import OverviewPage from './domains/dashboard/pages/OverviewPage.svelte';
  import LicensesPage from './domains/license/pages/LicensesPage.svelte';
  import ClientsPage from './domains/client/pages/ClientsPage.svelte';
  import LogsPage from './domains/log/pages/LogsPage.svelte';
  import SettingsPage from './domains/dashboard/pages/SettingsPage.svelte';

  const queryClient = new QueryClient();

  const routes = {
    '/': OverviewPage,
    '/licenses': LicensesPage,
    '/clients': ClientsPage,
    '/logs': LogsPage,
    '/settings': SettingsPage,
    '*': OverviewPage
  };

  // Sync hash to login if not logged in, or redirect to home if logged in and at /login
  $effect(() => {
    if (!$authStore.isLoggedIn) {
      if (window.location.hash !== '#/login') {
        window.location.hash = '#/login';
      }
    } else if (window.location.hash === '#/login') {
      window.location.hash = '#/';
    }
  });
</script>

<QueryClientProvider client={queryClient}>
  {#if $authStore.isLoggedIn}
    <DashboardLayout>
      <Router {routes} />
    </DashboardLayout>
  {:else}
    <Login />
  {/if}
</QueryClientProvider>
