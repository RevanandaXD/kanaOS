<script lang="ts">
	let today = $state(new Date());
	let view = $state(new Date(today.getFullYear(), today.getMonth(), 1));

	// refresh "today" every minute so day changes stay accurate
	$effect(() => {
		const id = setInterval(() => (today = new Date()), 60_000);
		return () => clearInterval(id);
	});

	const weekdays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

	const title = $derived(
		view.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
	);

	// 7-column grid, week starts Monday, includes days from previous/next months
	const cells = $derived.by(() => {
		const y = view.getFullYear();
		const m = view.getMonth();
		const offset = (new Date(y, m, 1).getDay() + 6) % 7;
		const daysInMonth = new Date(y, m + 1, 0).getDate();
		const total = Math.ceil((offset + daysInMonth) / 7) * 7;

		return Array.from({ length: total }, (_, i) => {
			const d = new Date(y, m, i - offset + 1);
			return { date: d, day: d.getDate(), current: d.getMonth() === m };
		});
	});

	const isToday = (d: Date) => d.toDateString() === today.toDateString();

	function shift(n: number) {
		view = new Date(view.getFullYear(), view.getMonth() + n, 1);
	}
	function reset() {
		view = new Date(today.getFullYear(), today.getMonth(), 1);
	}
</script>

<div class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 shadow-sm backdrop-blur-xl">
	<div class="mb-3 flex items-center justify-between">
		<button
			type="button"
			onclick={reset}
			title="Back to current month"
			class="text-sm font-medium capitalize text-(--text-primary) hover:text-(--accent)"
		>
			{title}
		</button>

		<div class="flex items-center gap-1">
			<button
				type="button"
				onclick={() => shift(-1)}
				aria-label="Previous month"
				class="flex h-7 w-7 items-center justify-center rounded-md text-(--text-secondary) hover:bg-(--surface-hover) hover:text-(--text-primary)"
			>
				<svg viewBox="0 0 24 24" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m15 18-6-6 6-6" /></svg>
			</button>
			<button
				type="button"
				onclick={() => shift(1)}
				aria-label="Next month"
				class="flex h-7 w-7 items-center justify-center rounded-md text-(--text-secondary) hover:bg-(--surface-hover) hover:text-(--text-primary)"
			>
				<svg viewBox="0 0 24 24" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m9 18 6-6-6-6" /></svg>
			</button>
		</div>
	</div>

	<div class="grid grid-cols-7 gap-y-1 text-center text-xs">
		{#each weekdays as w (w)}
			<span class="pb-1 text-[10px] tracking-wide text-(--text-muted) uppercase">{w}</span>
		{/each}

		{#each cells as cell (cell.date.getTime())}
			<span
				class="mx-auto flex h-7 w-7 items-center justify-center rounded-md tabular-nums
				{isToday(cell.date)
					? 'bg-(--accent-soft) font-medium text-(--accent) ring-1 ring-(--accent)/30'
					: cell.current
						? 'text-(--text-primary) hover:bg-(--surface-hover)'
						: 'text-(--text-muted) opacity-60'}"
			>
				{cell.day}
			</span>
		{/each}
	</div>
</div>