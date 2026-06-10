<script>
  import { createEventDispatcher } from 'svelte';
  import { X, Copy, CheckCircle2, Award, Calendar, RefreshCw } from 'lucide-svelte';

  const dispatch = createEventDispatcher();

  // Props
  export let clients = [];
  export let token = '';
  export let selectedClient = null; // can be pre-selected client object

  // State
  let clientID = selectedClient ? selectedClient.id : '';
  let trialLimit = 100;
  let useExpiration = false;
  let expirationDate = '';
  let generatedKey = '';
  let isSubmitting = false;
  let error = '';
  let copied = false;

  async function handleGenerate() {
    if (!clientID) {
      error = 'Silakan pilih nama klien toko.';
      return;
    }
    error = '';
    isSubmitting = true;

    try {
      let expiresAt = null;
      if (useExpiration && expirationDate) {
        expiresAt = new Date(expirationDate).toISOString();
      }

      const response = await fetch('/api/v1/admin/licenses', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          client_id: clientID,
          trial_limit: Number(trialLimit),
          expires_at: expiresAt
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Gagal membuat kunci lisensi baru.');
      }

      generatedKey = data.license_code;
      dispatch('success');
    } catch (err) {
      error = err.message;
    } finally {
      isSubmitting = false;
    }
  }

  async function copyKey() {
    if (!generatedKey) return;
    try {
      await navigator.clipboard.writeText(generatedKey);
      copied = true;
      setTimeout(() => copied = false, 2000);
    } catch (err) {
      console.error('Failed to copy', err);
    }
  }

  function handleClose() {
    dispatch('close');
  }
</script>

<div class="modal modal-open font-sans">
  <div class="modal-box bg-white border border-base-300 rounded-2xl shadow-2xl relative p-6 max-w-md">
    <!-- Close Button -->
    <button on:click={handleClose} class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4 text-gray-400 hover:text-gray-600">
      <X size={16} />
    </button>

    <div class="flex items-center gap-3 mb-6">
      <div class="p-2.5 bg-success/80 text-primary rounded-xl border border-success">
        <Award size={20} />
      </div>
      <h3 class="font-bold text-lg text-primary">Cetak Lisensi Toko</h3>
    </div>

    {#if error}
      <div class="alert alert-error mb-4 py-2 px-3 text-xs rounded-lg text-red-800 bg-error border border-red-200 font-semibold flex items-center gap-2">
        <span>{error}</span>
      </div>
    {/if}

    {#if !generatedKey}
      <!-- Form to Generate -->
      <div class="space-y-4">
        <!-- Client Selection -->
        <div class="form-control">
          <label class="label px-0.5 py-1" for="client-select">
            <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Pilih Klien Toko</span>
          </label>
          {#if selectedClient}
            <input 
              id="client-name"
              type="text" 
              value={selectedClient.name} 
              disabled 
              class="input input-bordered w-full bg-base-100 border-base-300 text-gray-500 rounded-xl cursor-not-allowed font-medium text-sm"
            />
          {:else}
            <select 
              id="client-select"
              bind:value={clientID}
              class="select select-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
            >
              <option value="" disabled selected>Pilih salah satu klien...</option>
              {#each clients as c}
                <option value={c.id}>{c.name}</option>
              {/each}
            </select>
          {/if}
        </div>

        <!-- Trial Limit -->
        <div class="form-control">
          <label class="label px-0.5 py-1" for="trial-input">
            <span class="label-text text-xs font-bold text-gray-500 uppercase tracking-wider">Batas Kuota Transaksi</span>
          </label>
          <div class="flex items-center gap-3">
            <input 
              id="trial-input"
              type="number" 
              min="0" 
              bind:value={trialLimit}
              class="input input-bordered w-28 bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm font-semibold"
            />
            <span class="text-[11px] text-gray-400 font-medium leading-relaxed">
              Isi <strong class="text-gray-700">0</strong> untuk versi tanpa batas (unlimited). Default 100.
            </span>
          </div>
        </div>

        <!-- Expiration Date Toggle -->
        <div class="form-control">
          <label class="label cursor-pointer justify-start gap-3 py-2 px-0.5" for="exp-toggle">
            <input 
              id="exp-toggle"
              type="checkbox" 
              bind:checked={useExpiration} 
              class="checkbox checkbox-primary checkbox-sm rounded-md" 
            />
            <span class="label-text text-xs font-bold text-gray-600 flex items-center gap-1.5 uppercase tracking-wide">
              <Calendar size={14} class="text-primary" />
              Tentukan Masa Kedaluwarsa
            </span>
          </label>
        </div>

        {#if useExpiration}
          <div class="form-control animate-fade-in">
            <input 
              type="datetime-local" 
              bind:value={expirationDate}
              class="input input-bordered w-full bg-white border-base-300 text-gray-800 rounded-xl focus:outline-none focus:border-primary text-sm"
            />
          </div>
        {/if}
      </div>

      <div class="modal-action gap-2 mt-6">
        <button on:click={handleClose} class="btn btn-outline btn-sm rounded-lg text-xs font-semibold h-9 px-4">Batal</button>
        <button 
          on:click={handleGenerate} 
          disabled={isSubmitting || (!selectedClient && !clientID)}
          class="btn btn-primary btn-sm text-white rounded-lg text-xs font-bold h-9 px-4"
        >
          {#if isSubmitting}
            <RefreshCw size={14} class="animate-spin mr-1 text-white" />
          {/if}
          Proses & Cetak Kunci
        </button>
      </div>
    {:else}
      <!-- Success display -->
      <div class="flex flex-col items-center justify-center py-4">
        <!-- Mint green circle check -->
        <div class="w-12 h-12 rounded-full bg-success/60 text-primary border border-success flex items-center justify-center mb-4">
          <CheckCircle2 size={24} />
        </div>
        <h4 class="font-bold text-primary text-md">Lisensi Berhasil Dicetak!</h4>
        <p class="text-xs text-gray-400 mt-1 text-center font-medium">Salin kode di bawah untuk diaktivasi klien secara daring.</p>
        
        <div class="w-full bg-base-100 border border-base-300 rounded-xl p-4 mt-6 flex items-center justify-between shadow-inner">
          <span class="font-mono text-lg font-bold text-primary tracking-wider select-all">{generatedKey}</span>
          <button 
            on:click={copyKey} 
            class="btn btn-circle btn-sm btn-ghost hover:bg-base-200 text-gray-500 hover:text-primary relative"
          >
            {#if copied}
              <CheckCircle2 size={16} class="text-primary" />
            {:else}
              <Copy size={16} />
            {/if}
          </button>
        </div>

        <!-- Details list in soft grey panel -->
        <div class="w-full mt-6 flex flex-col gap-2 bg-base-100 rounded-xl border border-base-300 p-4 text-xs text-gray-600 font-semibold leading-relaxed">
          <span class="font-bold text-primary uppercase text-[10px] tracking-widest block mb-1">Informasi Kunci:</span>
          <div>• Toko: <strong class="text-gray-800">{selectedClient ? selectedClient.name : clients.find(c => c.id === clientID)?.name}</strong></div>
          <div>• Kuota Transaksi: <strong class="text-gray-800">{trialLimit === 0 ? 'Tanpa Batas' : `${trialLimit} Transaksi`}</strong></div>
          {#if useExpiration && expirationDate}
            <div>• Masa Aktif: <strong class="text-gray-800">{new Date(expirationDate).toLocaleString('id-ID')}</strong></div>
          {/if}
        </div>

        <button on:click={handleClose} class="btn btn-primary btn-sm w-full mt-6 rounded-lg text-white font-bold h-10 text-xs">
          Selesai
        </button>
      </div>
    {/if}
  </div>
</div>
