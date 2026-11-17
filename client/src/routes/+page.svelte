<script lang="ts">
	import { onMount } from 'svelte';
	import { getServices } from '../api/services';
	import type { Category } from '../api/services';
	import Container from '../components/Container.svelte';

	let categories = $state<Category[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const data = await getServices();
			categories = data.categories ?? [];
		} catch (err) {
			error = 'Unable to reach the API';
			console.error(err);
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<p class="text-sm text-stone-500">Loading services…</p>
{:else if error}
	<div class="rounded-lg border border-stone-200 bg-white p-4 text-sm text-stone-700">
		<p class="font-medium text-stone-900">Could not load services</p>
		<p class="mt-1 text-stone-500">{error}</p>
	</div>
{:else if categories.length === 0}
	<p class="text-sm text-stone-500">No services configured.</p>
{:else}
	<div class="space-y-10">
		{#each categories as category (category.name)}
			<section>
				<div class="mb-4 flex items-baseline justify-between">
					<h2 class="text-sm font-medium tracking-wide text-stone-500 uppercase">
						{category.name}
					</h2>
					<span class="text-xs text-stone-400">{category.services.length}</span>
				</div>

				<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
					{#each category.services as service (service.url)}
						<Container {service} />
					{/each}
				</div>
			</section>
		{/each}
	</div>
{/if}