<script lang="ts">
	let today = $state(new Date());
	let view = $state(new Date(today.getFullYear(), today.getMonth(), 1));

	// perbarui "hari ini" tiap menit (supaya ganti hari tetap akurat)
	$effect(() => {
		const id = setInterval(() => (today = new Date()), 60_000);
		return () => clearInterval(id);
	});

	const weekdays = ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'];

	const title = $derived(
		view.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
	);

	// grid 7 kolom, minggu dimulai Senin, termasuk tanggal bulan sebelum/sesudahnya
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

<div class="rounded-lg border border-white/6 bg-white/4 p-4 backdrop-blur-lg">
	<div class="mb-3 flex items-center justify-between">
		<button
			type="button"
			onclick={reset}
			title="Kembali ke bulan ini"
			class="text-sm font-medium text-blue-100/90 capitalize hover:text-blue-50"
		>
			{title}
		</button>

		<div class="flex items-center gap-1">
			<button
				type="button"
				onclick={() => shift(-1)}
				aria-label="Bulan sebelumnya"
				class="flex h-6 w-6 items-center justify-center rounded-md text-blue-200/50 hover:bg-white/8 hover:text-blue-100"
			>
				<svg viewBox="0 0 24 24" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m15 18-6-6 6-6" /></svg>
			</button>
			<button
				type="button"
				onclick={() => shift(1)}
				aria-label="Bulan berikutnya"
				class="flex h-6 w-6 items-center justify-center rounded-md text-blue-200/50 hover:bg-white/8 hover:text-blue-100"
			>
				<svg viewBox="0 0 24 24" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m9 18 6-6-6-6" /></svg>
			</button>
		</div>
	</div>

	<div class="grid grid-cols-7 gap-y-1 text-center text-xs">
		{#each weekdays as w (w)}
			<span class="pb-1 text-[10px] tracking-wide text-blue-200/30 uppercase">{w}</span>
		{/each}

		{#each cells as cell (cell.date.getTime())}
			<span
				class="mx-auto flex h-7 w-7 items-center justify-center rounded-md tabular-nums
				{isToday(cell.date)
					? 'bg-blue-400/25 font-medium text-blue-50 ring-1 ring-blue-300/40'
					: cell.current
						? 'text-blue-100/70 hover:bg-white/8'
						: 'text-blue-200/20'}"
			>
				{cell.day}
			</span>
		{/each}
	</div>
</div>