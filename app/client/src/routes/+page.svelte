<script lang="ts">
	import { onMount } from 'svelte';
	// theme handled via CSS variables
	import { getContainers } from '../api/containers';
	import type { Container } from '../api/containers';
	import ContainerStatus from '../components/Container.svelte';
	import Clock from '../components/Clock.svelte';
	import Weather from '../components/Weather.svelte';
	import Calendar from '../components/Calendar.svelte';

	let containers = $state<Container[]>([]);
	let loading = $state(true);
	let error = $state('');

	const running = $derived(containers.filter((c) => c.state === 'running').length);
	const stopped = $derived(
		containers.filter((c) => c.state === 'exited' || c.state === 'dead').length
	);
	const otherState = $derived(
		containers.filter((c) => !['running', 'exited', 'dead'].includes(c.state)).length
	);
	const issues = $derived(
		containers.filter((c) => c.state !== 'running' || c.health === 'unhealthy')
	);

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

				// Keep previously loaded data visible; only update the error state.
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

<div
	class="grid grid-cols-1 gap-4 sm:gap-6 lg:grid-cols-[240px_minmax(0,1fr)_240px] xl:grid-cols-[260px_minmax(0,1fr)_260px]"
>
	<!-- Left -->
	<aside class="space-y-3 sm:space-y-4 lg:sticky lg:top-6 lg:self-start">
		<Clock />
		<Calendar />
		<Weather />
	</aside>

	<!-- Center -->
	<main>
		<section>
			<div
				class="mb-3 flex items-center justify-between rounded-full border border-(--border-subtle) bg-(--surface-card) px-4 py-2 shadow-sm backdrop-blur-xl transition-colors sm:mb-4 sm:px-5"
			>
				<h2 class="text-lg font-semibold tracking-wide text-(--text-primary) sm:text-xl">
					Containers
				</h2>
				<span class="text-xs text-(--text-secondary) tabular-nums">
					{running}/{containers.length} running
				</span>
			</div>

			{#if loading}
				<p class="text-sm text-(--text-muted)">Reading Docker…</p>
			{:else if error && containers.length === 0}
				<div
					class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 text-sm shadow-sm backdrop-blur-xl"
				>
					<p class="font-medium text-(--text-primary)">Docker unavailable</p>
					<p class="mt-1 text-(--text-secondary)">{error}</p>
				</div>
			{:else if containers.length === 0}
				<p class="text-sm text-(--text-muted)">No containers found.</p>
			{:else}
				{#if error}
					<p class="mb-3 text-xs text-amber-400/80">
						Last refresh failed ({error}). Showing the last known state.
					</p>
				{/if}

				<div class="grid grid-cols-2 gap-2 sm:grid-cols-2 sm:gap-2.5 xl:grid-cols-3">
					{#each containers as container (container.id)}
						<ContainerStatus {container} />
					{/each}
				</div>
			{/if}
		</section>
	</main>

	<!-- Right -->
	<aside class="space-y-3 sm:space-y-4 lg:sticky lg:top-6 lg:self-start">
		<div
			class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 shadow-sm backdrop-blur-xl sm:p-5"
		>
			<p class="text-xs font-medium tracking-wide text-(--text-muted) uppercase">Summary</p>
			<div class="mt-3 grid grid-cols-2 gap-2">
				<div class="rounded-md bg-(--surface-inset) px-3 py-2">
					<span class="block text-2xl font-light text-(--text-primary)">{running}</span>
					<span class="text-xs text-(--text-secondary)">Running</span>
				</div>
				<div class="rounded-md bg-(--surface-inset) px-3 py-2">
					<span class="block text-2xl font-light text-(--text-primary)">{containers.length}</span>
					<span class="text-xs text-(--text-secondary)">Total</span>
				</div>
			</div>
		</div>

		<div
			class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 shadow-sm backdrop-blur-xl sm:p-5"
		>
			<p class="text-xs font-medium tracking-wide text-(--text-muted) uppercase">
				Container status
			</p>

			<div class="mt-3 grid grid-cols-3 gap-2 text-center">
				<div class="rounded-md bg-(--surface-inset) px-2 py-2">
					<span class="block text-lg font-light text-emerald-400 tabular-nums">{running}</span>
					<span class="text-[11px] text-(--text-secondary)">Running</span>
				</div>
				<div class="rounded-md bg-(--surface-inset) px-2 py-2">
					<span class="block text-lg font-light text-(--text-primary) tabular-nums">{stopped}</span>
					<span class="text-[11px] text-(--text-secondary)">Stopped</span>
				</div>
				<div class="rounded-md bg-(--surface-inset) px-2 py-2">
					<span class="block text-lg font-light text-amber-400 tabular-nums">{otherState}</span>
					<span class="text-[11px] text-(--text-secondary)">Other</span>
				</div>
			</div>

			{#if loading}
				<p class="mt-4 text-sm text-(--text-muted)">Checking status…</p>
			{:else if issues.length === 0}
				<p class="mt-4 text-sm text-emerald-400/90">All containers are running normally.</p>
			{:else}
				<div class="mt-4 border-t border-(--border-subtle) pt-3">
					<p class="mb-2 text-[11px] font-medium tracking-wide text-amber-400 uppercase">
						Needs attention ({issues.length})
					</p>
					<ul class="space-y-2">
						{#each issues as item (item.id)}
							<li class="min-w-0">
								<div class="flex items-center gap-2">
									<span
										class="h-1.5 w-1.5 shrink-0 rounded-full {item.health === 'unhealthy'
											? 'bg-rose-400'
											: item.state === 'running'
												? 'bg-amber-400'
												: 'bg-(--text-muted)'}"
										aria-hidden="true"
									></span>
									<span class="truncate text-sm text-(--text-primary)">
										{item.name || item.id.slice(0, 12)}
									</span>
								</div>
								<p
									class="mt-0.5 truncate pl-3.5 text-xs text-(--text-secondary)"
									title={item.status}
								>
									{item.health === 'unhealthy' ? 'Unhealthy · ' : ''}{item.status}
								</p>
							</li>
						{/each}
					</ul>
				</div>
			{/if}
		</div>
	</aside>
</div>
