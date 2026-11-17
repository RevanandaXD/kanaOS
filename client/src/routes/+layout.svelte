<script lang="ts">
	import '@fontsource-variable/roboto'
	import '@fontsource-variable/fira-code';

	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import Navbar from '../components/Navbar.svelte';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>KanaOS</title>
</svelte:head>

<div class="min-h-screen bg-[#0a0e1a] text-blue-100">
	<div class="px-4 pt-4 sm:px-6">
		<Navbar />
	</div>
	<main class="max-w-7xl mx-auto py-7.5">
		{@render children()}
	</main>
</div>
