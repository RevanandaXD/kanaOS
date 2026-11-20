<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import Dropdown from './Dropdown.svelte';
	import { themeStore } from '#lib/theme.svelte.ts';
	import { getSystem } from '../api/system';

	let isOpen = $state(false);
	let menuRef = $state<HTMLElement | null>(null);
	let osName = $state('');

	const isLight = $derived(themeStore.current === 'light');

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
	class="relative z-50 mx-auto flex max-w-7xl items-center justify-between rounded-full border px-6 py-1 backdrop-blur-lg transition-colors {isLight
		? 'border-slate-200/70 bg-white/70 text-slate-700'
		: 'border-white/6 bg-white/4 text-blue-100'}"
>
	<a href="/" class="text-sm font-medium tracking-tight {isLight ? 'text-slate-900' : 'text-blue-100/90'}">KanaOS</a>

	<nav bind:this={menuRef} class="relative hidden items-center gap-5 text-xs sm:flex">
		<button
			type="button"
			class="transition-colors {isLight
				? isOpen
					? 'text-slate-900'
					: 'text-slate-500 hover:text-slate-800'
				: isOpen
					? 'text-blue-100/90'
					: 'text-blue-200/40 hover:text-blue-100/80'}"
			onclick={() => (isOpen = !isOpen)}
		>
			Dashboard
		</button>
		{#if isOpen}
			<div
				class="absolute left-1/2 top-full mt-3 -translate-x-1/2"
				transition:fly={{ y: -12, duration: 180 }}
			>
				<Dropdown />
			</div>
		{/if}
	</nav>

	<div class="flex items-center gap-3 text-xs tabular-nums {isLight ? 'text-slate-500' : 'text-blue-200/50'}">
		{#if osName}
			<p title={osName}>{osName}</p>
		{/if}
	</div>
</header>
