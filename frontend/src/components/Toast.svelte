<script>
    import { CheckCircle2, Info } from 'lucide-svelte';

    export let message = '';
    export let type = 'success'; // 'success' | 'error' | 'info'
</script>

{#if message}
    <div class="toast toast-end toast-bottom z-50 animate-fade-in">
        <div class="alert shadow-2xl py-3 px-5 rounded-md text-xs border border-base-300 flex items-center gap-2
          {type === 'success' ? 'bg-success text-primary border-success' : ''}
          {type === 'error' ? 'bg-error text-red-800 border-red-200' : ''}
          {type === 'info' ? 'bg-blue-50 text-blue-800 border-blue-100' : ''}
        ">
            {#if type === 'success'}
                <CheckCircle2 size={14} class="text-primary"/>
            {:else}
                <Info size={14}/>
            {/if}
            <span class="font-bold">{message}</span>
        </div>
    </div>
{/if}
