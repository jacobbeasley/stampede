<script>
  import KitchenSink from "./KitchenSink.svelte";

  let name = $state("World");
  let showKitchenSink = $state(false);
</script>

<div class="hero min-h-[40vh] bg-base-200 rounded-2xl">
  <div class="hero-content text-center flex-col gap-6">
    <div class="max-w-lg">
      <h2
        class="text-5xl font-extrabold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent"
      >
        Hello! 👋
      </h2>
      <p class="py-4 text-base-content/70 text-lg">
        A starter app for efficient development of multi-tenant SaaS applications.
      </p>
      <p class="py-4 text-base-content/70 text-lg">
        Rapid. Efficient. Scaleable. 
      </p>
      <p class="py-4 text-base-content/70 text-lg">
        Buffalo · Svelte · Tailwind CSS · DaisyUI — all wired up and ready to
        go.
      </p>
    </div>

    <div class="flex flex-wrap gap-2 justify-center mt-2">
      <span class="badge badge-primary badge-outline">Buffalo v0.18</span>
      <span class="badge badge-secondary badge-outline">Svelte 5</span>
      <span class="badge badge-accent badge-outline">Tailwind v4</span>
      <span class="badge badge-info badge-outline">DaisyUI</span>
    </div>
    <div class="mt-4">
      <button 
        class="btn btn-outline" 
        onclick={() => showKitchenSink = !showKitchenSink}
      >
        {showKitchenSink ? "Hide Kitchen Sink" : "Show Kitchen Sink"}
      </button>
    </div>
  </div>
</div>

{#if showKitchenSink}
  <KitchenSink />
{/if}
