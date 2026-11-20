<script lang="ts">
	let now = $state(new Date());

	$effect(() => {
		const id = setInterval(() => (now = new Date()), 1000);
		return () => clearInterval(id);
	});

	const time = $derived(
		now.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })
	);
	const seconds = $derived(String(now.getSeconds()).padStart(2, '0'));
	const date = $derived(
		now.toLocaleDateString('en-US', {
			weekday: 'long',
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		})
	);
</script>

<div
	class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 shadow-sm backdrop-blur-xl sm:p-5"
>
	<div
		class="flex items-baseline gap-1 font-light tracking-tight text-(--text-primary) tabular-nums"
	>
		<span class="text-4xl sm:text-5xl">{time}</span>
		<span class="text-base text-(--text-muted) sm:text-lg">{seconds}</span>
	</div>
	<p class="mt-2 text-xs text-(--text-secondary) sm:text-sm">{date}</p>
</div>
