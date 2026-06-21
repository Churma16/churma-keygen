<script>
  export let type = 'button';
  export let variant = 'primary'; // 'primary' | 'secondary' | 'outline' | 'ghost' | 'error' | 'success'
  export let size = 'md'; // 'sm' | 'md' | 'lg'
  export let disabled = false;
  export let className = '';

  const variants = {
    primary: 'btn-primary text-white shadow-md shadow-primary/10 hover:shadow-primary/25',
    secondary: 'btn-secondary text-white shadow-md shadow-secondary/10 hover:shadow-secondary/25',
    outline: 'btn-outline border-base-300',
    ghost: 'btn-ghost hover:bg-base-200 text-gray-500 hover:text-primary',
    error: 'btn-error text-white',
    success: 'btn-success text-white'
  };

  const sizes = {
    sm: 'btn-sm text-xs h-9 px-4',
    md: 'h-12 text-sm px-6',
    lg: 'btn-lg text-base h-14 px-8'
  };
</script>

<button
  {type}
  disabled={disabled}
  on:click
  class="btn rounded-md font-bold tracking-wide active:scale-[0.99] transition duration-150 {variants[variant] || ''} {sizes[size] || ''} {className}"
>
  <slot />
</button>
