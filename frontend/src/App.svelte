<script>
  import Login from './domains/auth/pages/Login.svelte';
  import Dashboard from './domains/dashboard/pages/Dashboard.svelte';
  import { authStore } from './domains/auth/store/authStore';
</script>

{#if $authStore.isLoggedIn}
  <Dashboard />
{:else}
  <Login />
{/if}
