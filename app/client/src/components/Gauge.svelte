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
	class="flex flex-col items-center gap-2 rounded-md border border-(--border-subtle) bg-(--surface-inset) p-3 transition hover:border-(--text-muted)/30"
>
	<div class="relative h-24 w-24">
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
			<span class="font-mono text-base font-semibold text-(--text-primary) tabular-nums">
				{Math.round(value)}%
			</span>
		</div>
	</div>
	<div class="text-center">
		<p class="text-xs font-medium text-(--text-primary)">{label}</p>
		<p class="mt-0.5 text-[10px] text-(--text-muted)">{sub}</p>
	</div>
</div>