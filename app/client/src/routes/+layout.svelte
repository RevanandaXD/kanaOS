<script lang="ts">
	import '@fontsource-variable/roboto';
	import '@fontsource-variable/fira-code';

	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import Navbar from '../components/Navbar.svelte';
	import { wallpaperStore } from '#lib/wallpaper.svelte.ts';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();

	const tiles = Array.from({ length: 48 }, (_, i) => i);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>KanaOS</title>
</svelte:head>

<div class="app-root relative min-h-screen overflow-hidden text-(--text-primary) transition-colors">
	<div class="pointer-events-none fixed inset-0 z-0">
		<div
			class="absolute inset-0 bg-cover bg-center bg-no-repeat"
			style="background-image: var(--kanaos-wallpaper, url('/wallpaper-default1.jpg'))"
		></div>

		{#if wallpaperStore.previous}
			{#key `${wallpaperStore.previous}-${wallpaperStore.current}-${wallpaperStore.transition}`}
				{#if wallpaperStore.transition === 'fade'}
					<div
						class="wallpaper-fade absolute inset-0 bg-cover bg-center bg-no-repeat"
						style={`background-image: url('${wallpaperStore.previous}')`}
					></div>
				{:else if wallpaperStore.transition === 'wave'}
					<div
						class="wallpaper-wave absolute inset-0 bg-cover bg-center bg-no-repeat"
						style={`background-image: url('${wallpaperStore.previous}')`}
					></div>
				{:else if wallpaperStore.transition === 'tiles'}
					<div class="absolute inset-0 grid grid-cols-8 grid-rows-6">
						{#each tiles as i}
							<div
								class="wallpaper-tile bg-cover bg-center bg-no-repeat"
								style={`background-image: url('${wallpaperStore.previous}'); background-size: 800% 600%; background-position: ${(i % 8) * (100 / 7)}% ${Math.floor(i / 8) * (100 / 5)}%; animation-delay: ${(i % 8) * 40 + Math.floor(i / 8) * 35}ms;`}
							></div>
						{/each}
					</div>
				{/if}
			{/key}
		{/if}

		<div class="theme-overlay absolute inset-0"></div>
	</div>

	<div class="relative z-10">
		<div class="px-3 pt-3 sm:px-6 sm:pt-4">
			<Navbar />
		</div>
		<main class="mx-auto max-w-7xl px-3 py-4 sm:px-6 sm:py-6 lg:py-7">
			{@render children()}
		</main>
	</div>
</div>

<style>
	.wallpaper-fade {
		animation: wallpaper-fade 700ms ease forwards;
	}
	@keyframes wallpaper-fade {
		from {
			opacity: 1;
		}
		to {
			opacity: 0;
		}
	}

	.wallpaper-wave {
		animation: wallpaper-wave 700ms ease-in-out forwards;
	}
	@keyframes wallpaper-wave {
		0% {
			clip-path: polygon(0 0, 100% 0, 100% 100%, 0 100%);
			filter: saturate(1);
		}
		30% {
			clip-path: polygon(0 0, 82% 0, 68% 20%, 80% 40%, 64% 60%, 78% 80%, 62% 100%, 0 100%);
		}
		60% {
			clip-path: polygon(0 0, 42% 0, 28% 20%, 40% 40%, 24% 60%, 36% 80%, 20% 100%, 0 100%);
		}
		100% {
			clip-path: polygon(0 0, 0 0, 0 20%, 0 40%, 0 60%, 0 80%, 0 100%, 0 100%);
			opacity: 0.2;
			filter: saturate(1.1);
		}
	}

	.wallpaper-tile {
		animation: wallpaper-tile 700ms ease forwards;
	}
	@keyframes wallpaper-tile {
		0% {
			opacity: 1;
			transform: scale(1);
		}
		100% {
			opacity: 0;
			transform: scale(0.2) rotate(8deg);
		}
	}
</style>
