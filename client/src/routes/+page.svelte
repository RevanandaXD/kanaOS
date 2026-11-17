<script lang="ts">
	import { onMount } from 'svelte';
	import { getServices } from '../api/services';
	import type { Category } from '../api/services';
	import Container from '../components/Container.svelte';
	import Clock from '../components/Clock.svelte';
	import Weather from '../components/Weather.svelte';
	import Calendar from '../components/Calendar.svelte';

	let categories = $state<Category[]>([]);
	let loading = $state(true);
	let error = $state('');

	const totalServices = $derived(categories.reduce((n, c) => n + c.services.length, 0));

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

<div class="grid grid-cols-1 gap-6 lg:grid-cols-[260px_minmax(0,1fr)_260px]">
	<!-- Kiri -->
	<aside class="space-y-4 lg:sticky lg:top-6 lg:self-start">
		<Clock />
		<Calendar />
		<Weather />
	</aside>

	<!-- Tengah -->
	<main>
		{#if loading}
			<p class="text-sm text-blue-200/35">Loading services…</p>
		{:else if error}
			<div class="rounded-lg border border-white/6 bg-white/4 p-4 text-sm backdrop-blur-lg">
				<p class="font-medium text-blue-100/90">Could not load services</p>
				<p class="mt-1 text-blue-200/40">{error}</p>
			</div>
		{:else if categories.length === 0}
			<p class="text-sm text-blue-200/35">No services configured.</p>
		{:else}
			<div class="space-y-10">
				{#each categories as category (category.name)}
					<section>
						<div class="mb-4 flex items-baseline justify-between">
							<h2 class="text-xl font-semibold tracking-wide text-blue-200/40">
								Infrastructure
							</h2>
						</div>
						<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
							{#each category.services as service (service.url)}
								<Container {service} />
							{/each}
						</div>
					</section>
				{/each}
			</div>
		{/if}
	</main>

	<!-- Kanan -->
	<aside class="space-y-4 lg:sticky lg:top-6 lg:self-start">
		<div class="rounded-lg border border-white/[0.06] bg-white/[0.04] p-5 backdrop-blur-lg">
			<p class="text-xs font-medium tracking-wide text-blue-200/40 uppercase">Ringkasan</p>
			<div class="mt-3 grid grid-cols-2 gap-2">
				<div class="rounded-md bg-white/[0.03] px-3 py-2">
					<span class="block text-2xl font-light text-blue-100/90">{categories.length}</span>
					<span class="text-xs text-blue-200/35">Kategori</span>
				</div>
				<div class="rounded-md bg-white/[0.03] px-3 py-2">
					<span class="block text-2xl font-light text-blue-100/90">{totalServices}</span>
					<span class="text-xs text-blue-200/35">Service</span>
				</div>
			</div>
		</div>

		<div class="rounded-lg border border-dashed border-white/[0.08] p-5 text-sm text-blue-200/20">
			Slot widget kosong
		</div>
	</aside>
</div>
