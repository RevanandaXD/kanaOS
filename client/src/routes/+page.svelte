<script lang="ts">
	import { onMount } from 'svelte';
	import { themeStore } from '#lib/theme.svelte.ts';
	import { getContainers } from '../api/containers';
	import type { Container } from '../api/containers';
	import ContainerStatus from '../components/Container.svelte';
	import Clock from '../components/Clock.svelte';
	import Weather from '../components/Weather.svelte';
	import Calendar from '../components/Calendar.svelte';

	let containers = $state<Container[]>([]);
	let loading = $state(true);
	let error = $state('');

	const isLight = $derived(themeStore.current === 'light');
	const running = $derived(containers.filter((c) => c.state === 'running').length);

	onMount(() => {
		let active = true;

		const load = async () => {
			try {
				const data = await getContainers();
				if (!active) return;

				containers = data.containers ?? [];
				error = '';
			} catch (err) {
				if (!active) return;

				// Data lama sengaja dibiarkan tampil; cuma status error yang berubah.
				error = err instanceof Error ? err.message : 'Unable to read Docker';
			} finally {
				if (active) loading = false;
			}
		};

		load();
		const interval = setInterval(load, 5000);

		return () => {
			active = false;
			clearInterval(interval);
		};
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
		<section>
			<div
				class="mb-4 flex items-center justify-between rounded-full border px-5 py-2 backdrop-blur-lg transition-colors {isLight
					? 'border-slate-200/70 bg-white/70 text-slate-700'
					: 'border-white/6 bg-white/4 text-blue-100'}"
			>
				<h2
					class="text-xl font-semibold tracking-wide {isLight
						? 'text-slate-900'
						: 'text-blue-200/70'}"
				>
					Containers
				</h2>
				<span class="text-xs tabular-nums {isLight ? 'text-slate-500' : 'text-blue-200/30'}">
					{running}/{containers.length} running
				</span>
			</div>

			{#if loading}
				<p class="text-sm text-blue-200/35">Reading Docker…</p>
			{:else if error && containers.length === 0}
				<div class="rounded-lg border border-white/6 bg-white/4 p-4 text-sm backdrop-blur-lg">
					<p class="font-medium text-blue-100/90">Docker unavailable</p>
					<p class="mt-1 text-blue-200/40">{error}</p>
				</div>
			{:else if containers.length === 0}
				<p class="text-sm text-blue-200/35">No containers found.</p>
			{:else}
				{#if error}
					<p class="mb-3 text-xs text-amber-300/60">
						Last refresh failed ({error}). Showing the last known state.
					</p>
				{/if}

				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
					{#each containers as container (container.id)}
						<ContainerStatus {container} />
					{/each}
				</div>
			{/if}
		</section>
	</main>

	<!-- Kanan -->
	<aside class="space-y-4 lg:sticky lg:top-6 lg:self-start">
		<div class="rounded-lg border border-white/6 bg-white/4 p-5 backdrop-blur-lg">
			<p class="text-xs font-medium tracking-wide text-blue-200/40 uppercase">Ringkasan</p>
			<div class="mt-3 grid grid-cols-2 gap-2">
				<div class="rounded-md bg-white/3 px-3 py-2">
					<span class="block text-2xl font-light text-blue-100/90">{running}</span>
					<span class="text-xs text-blue-200/35">Running</span>
				</div>
				<div class="rounded-md bg-white/3 px-3 py-2">
					<span class="block text-2xl font-light text-blue-100/90">{containers.length}</span>
					<span class="text-xs text-blue-200/35">Total</span>
				</div>
			</div>
		</div>

		<div class="rounded-lg border border-dashed border-white/8 p-5 text-sm text-blue-200/20">
			Slot widget kosong
		</div>
	</aside>
</div>
