<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import Dropdown from './Dropdown.svelte';
	import { getSystem } from '../api/system';

	let isOpen = $state(false);
	let menuRef = $state<HTMLElement | null>(null);
	let osName = $state('');

	function handleWindowClick(event: MouseEvent) {
		if (isOpen && menuRef && !menuRef.contains(event.target as Node)) {
			isOpen = false;
		}
	}

	onMount(async () => {
		try {
			const info = await getSystem();
			osName = info.os || info.hostname || '';
		} catch {
			// biarkan kosong jika info sistem belum tersedia
		}
	});
</script>

<svelte:window onclick={handleWindowClick} />

<header
	class="relative z-50 mx-auto flex max-w-7xl items-center justify-between rounded-full border border-(--border-subtle) bg-(--surface-nav) px-4 py-1.5 text-(--text-primary) shadow-sm backdrop-blur-xl transition-colors sm:px-6"
>
	<a href="/" class="text-sm font-semibold tracking-tight text-(--text-primary)">KanaOS</a>

	<nav bind:this={menuRef} class="relative flex items-center gap-4 text-xs sm:gap-5">
		<button
			type="button"
			class="rounded-md px-2 py-1 transition-colors {isOpen
				? 'text-(--text-primary)'
				: 'text-(--text-secondary) hover:text-(--text-primary)'}"
			onclick={() => (isOpen = !isOpen)}
		>
			Dashboard
		</button>
		{#if isOpen}
			<div
				class="fixed top-16 left-1/2 z-50 w-[min(92vw,45rem)] -translate-x-1/2 sm:absolute sm:top-full sm:left-1/2 sm:mt-3 sm:w-auto sm:min-w-[45rem]"
				transition:fly={{ y: -12, duration: 180 }}
			>
				<Dropdown />
			</div>
		{/if}
	</nav>

	<div
		class="hidden max-w-[30%] min-w-0 items-center justify-end text-xs text-(--text-muted) tabular-nums sm:flex"
	>
		{#if osName}
			<p class="truncate" title={osName}>{osName}</p>
		{/if}
	</div>
</header>
