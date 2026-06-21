<script>
  export let id = '';
  export let type = 'text';
  export let placeholder = '';
  export let value = '';
  export let required = false;
  export let disabled = false;
  export let label = '';
  export let className = '';
</script>

<div class="form-control w-full {className}">
  {#if label}
    <label class="label px-0.5 py-1" for={id}>
      <span class="label-text font-bold text-gray-700 text-xs uppercase tracking-wider">{label}</span>
    </label>
  {/if}
  <div class="relative">
    {#if $$slots.icon}
      <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
        <slot name="icon" />
      </div>
    {/if}
    <input
      {id}
      {type}
      {placeholder}
      bind:value
      {required}
      {disabled}
      class="input input-bordered w-full bg-gray-50 border-gray-300 text-gray-800 rounded-md focus:bg-white focus:outline-none focus:border-primary transition duration-150 text-sm {$$slots.icon ? 'pl-10' : ''}"
      on:input
    />
  </div>
</div>
