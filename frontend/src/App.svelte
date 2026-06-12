<script>
  import Login from './pages/Login.svelte';
  import Dashboard from './pages/Dashboard.svelte';
  import { authStore } from './stores/authStore';
</script>

{#if $authStore.isLoggedIn}
  <Dashboard />
{:else}
  <Login />
{/if}
