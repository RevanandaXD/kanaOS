<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import Dropdown from './Dropdown.svelte';
	import { getSystem } from '../api/system';
	import { version } from '../../package.json' with { type: 'json' };

	let isOpen = $state(false);
	let headerRef = $state<HTMLElement | null>(null);
	let dropdownRef = $state<HTMLElement | null>(null);
	let osName = $state('');

	function handleWindowPointerDown(event: PointerEvent) {
		if (!isOpen) return;
		const path = event.composedPath();
		const insideHeader = headerRef ? path.includes(headerRef) : false;
		const insideDropdown = dropdownRef ? path.includes(dropdownRef) : false;
		if (!insideHeader && !insideDropdown) {
			isOpen = false;
		}
	}

	function handleWindowKeyDown(event: KeyboardEvent) {
		if (isOpen && event.key === 'Escape') {
			isOpen = false;
		}
	}

	onMount(async () => {
		try {
			const info = await getSystem();
			osName = info.os || info.hostname || '';
		} catch {

		}
	});
</script>

<svelte:window onpointerdown={handleWindowPointerDown} onkeydown={handleWindowKeyDown} />

<header
	bind:this={headerRef}
	class="relative z-50 mx-auto flex max-w-7xl items-center justify-between rounded-full border border-(--border-subtle) bg-(--surface-nav) px-4 py-2 text-(--text-primary) shadow-sm backdrop-blur-xl transition-colors sm:px-6 sm:py-1.5"
>
	<a href="/" class="shrink-0 text-sm font-semibold tracking-tight text-(--text-primary)">
		KanaOS
	</a>

	<nav class="relative flex items-center justify-between">
		<button
			type="button"
			aria-expanded={isOpen}
			class="inline-flex h-8 items-center gap-1.5 rounded-full px-3 text-xs font-medium transition-colors {isOpen
				? 'bg-(--surface-hover) text-(--text-primary)'
				: 'text-(--text-secondary) hover:bg-(--surface-hover) hover:text-(--text-primary)'}"
			onclick={() => (isOpen = !isOpen)}
		>
			<span>Dashboard</span>
		</button>

		{#if isOpen}
			<div
				class="absolute top-full left-1/2 z-50 mt-3 hidden w-[min(92vw,45rem)] -translate-x-1/2 sm:block"
				transition:fly={{ y: -12, duration: 180 }}
			>
				<Dropdown />
			</div>
		{/if}

	</nav>

	<div class="hidden max-w-[30%] min-w-0 items-center justify-end text-xs text-(--text-muted) tabular-nums sm:flex">
		{version}
	</div>
</header>

{#if isOpen}
	<div
		bind:this={dropdownRef}
		class="fixed inset-x-3 top-16 z-40 sm:hidden"
		transition:fly={{ y: -12, duration: 180 }}
	>
		<Dropdown />
	</div>
{/if}
