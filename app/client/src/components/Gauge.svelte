<script lang="ts">
	import type { Component } from 'svelte';

	let {
		value = 0,
		label,
		sub = '',
		color = 'var(--accent)',
		Icon
	}: {
		value?: number;
		label: string;
		sub?: string;
		color?: string;
		Icon: Component<any>;
	} = $props();

	const r = 34;
	const c = 2 * Math.PI * r;
	const offset = $derived(c * (1 - Math.min(100, Math.max(0, value)) / 100));
</script>

<div
	class="flex flex-col items-center gap-1.5 rounded-md border border-(--border-subtle) bg-(--surface-inset) p-2 transition hover:border-(--text-muted)/30 sm:gap-2 sm:p-3"
>
	<div class="relative h-16 w-16 sm:h-24 sm:w-24">
		<svg viewBox="0 0 80 80" class="h-full w-full -rotate-90">
			<circle cx="40" cy="40" {r} fill="none" stroke="var(--surface-hover)" stroke-width="6" />
			<circle
				cx="40"
				cy="40"
				{r}
				fill="none"
				stroke={color}
				stroke-width="6"
				stroke-linecap="round"
				stroke-dasharray={c}
				stroke-dashoffset={offset}
				class="transition-[stroke-dashoffset] duration-700 ease-out"
			/>
		</svg>
		<div class="absolute inset-0 flex flex-col items-center justify-center gap-0.5">
			<span style="color:{color}"><Icon size={14} weight="bold" /></span>
			<span class="font-mono text-xs font-semibold text-(--text-primary) tabular-nums sm:text-base">
				{Math.round(value)}%
			</span>
		</div>
	</div>
	<div class="min-w-0 max-w-full text-center">
		<p class="text-[11px] font-medium text-(--text-primary) sm:text-xs">{label}</p>
		<p class="mt-0.5 truncate text-[9px] text-(--text-muted) sm:text-[10px]">{sub}</p>
	</div>
</div>