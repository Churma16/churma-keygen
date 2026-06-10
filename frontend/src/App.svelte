<script>
  import Login from './components/Login.svelte';
  import Dashboard from './components/Dashboard.svelte';
  import { onMount } from 'svelte';

  let isLoggedIn = false;

  onMount(() => {
    const token = localStorage.getItem('admin_token');
    isLoggedIn = !!token;
  });

  function handleLoginSuccess() {
    isLoggedIn = true;
  }
</script>

{#if isLoggedIn}
  <Dashboard />
{:else}
  <Login on:loginSuccess={handleLoginSuccess} />
{/if}
