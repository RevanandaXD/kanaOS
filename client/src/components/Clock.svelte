<script lang="ts">
	let now = $state(new Date());

	$effect(() => {
		const id = setInterval(() => (now = new Date()), 1000);
		return () => clearInterval(id);
	});

	const time = $derived(
		now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', hour12: false })
	);
	const seconds = $derived(String(now.getSeconds()).padStart(2, '0'));
	const date = $derived(
		now.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
	);
</script>

<div class="rounded-lg border border-white/6 bg-white/4 p-5 backdrop-blur-lg">
	<div class="flex items-baseline gap-1 font-light tracking-tight text-blue-100/90 tabular-nums">
		<span class="text-5xl">{time}</span>
		<span class="text-lg text-blue-200/30">{seconds}</span>
	</div>
	<p class="mt-2 text-sm text-blue-200/40">{date}</p>
</div>
