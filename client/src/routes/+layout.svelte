<script lang="ts">
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

<div class="min-h-screen bg-stone-50 text-stone-900">
	<Navbar />
	<main class="mx-auto max-w-5xl px-6 py-10">
		{@render children()}
	</main>
</div>
