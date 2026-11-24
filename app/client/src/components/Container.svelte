<script lang="ts">
	import CubeIcon from 'phosphor-svelte/lib/Cube';
	import { dockerIcon } from '../lib/docker-icon';
	import type { Container, Port } from '../api/containers';

	let { container }: { container: Container } = $props();
	let iconFailed = $state(false);

	const title = $derived(container.name || container.id.slice(0, 12));
	const icon = $derived(dockerIcon(container.image));
	const link = $derived(containerLink(container));

	const tone = $derived.by(() => {
		if (container.state !== 'running') return 'bg-(--text-muted)';
		if (container.health === 'unhealthy') return 'bg-rose-400';
		if (container.health === 'starting') return 'bg-amber-400';
		return 'bg-emerald-400';
	});

	function portLabel(port: Port): string {
		return port.public
			? `${port.public}→${port.private}/${port.type}`
			: `${port.private}/${port.type}`;
	}

	function containerLink(item: Container): string | null {
		if (typeof window === 'undefined') return null;

		const port = item.ports.find((p) => p.public);
		if (!port?.public) return null;

		const scheme = port.public === 443 ? 'https' : 'http';

		return `${scheme}://${window.location.hostname}:${port.public}`;
	}
</script>

<a
	href={link ?? undefined}
	target={link ? '_blank' : undefined}
	rel="noreferrer"
	class="group flex flex-col gap-2 rounded-xl border border-(--border-subtle) bg-(--surface-card) p-2.5 shadow-sm backdrop-blur-xl transition sm:p-3 {link
		? 'hover:-translate-y-0.5 hover:bg-(--surface-card-hover) hover:border-(--border-subtle)'
		: 'cursor-default'}"
>
	<div class="flex items-start justify-between gap-2">
		{#if icon && !iconFailed}
			<img
				src={icon}
				alt=""
				class="h-8 w-8 rounded-md object-contain sm:h-9 sm:w-9"
				loading="lazy"
				onerror={() => (iconFailed = true)}
			/>
		{:else}
			<div
				class="flex h-8 w-8 items-center justify-center rounded-md bg-(--surface-inset) text-(--text-muted) sm:h-9 sm:w-9"
				aria-hidden="true"
			>
				<CubeIcon size={18} weight="duotone" />
			</div>
		{/if}

		<span
			class="flex items-center gap-1 rounded-full bg-(--surface-inset) px-1.5 py-0.5 text-[10px] text-(--text-secondary)"
		>
			<span class="h-1.5 w-1.5 rounded-full {tone}" aria-hidden="true"></span>
			{container.state}
		</span>
	</div>

	<div class="min-w-0 flex-1">
		<h3 class="truncate text-xs font-semibold text-(--text-primary) sm:text-[13px]">{title}</h3>
		<p
			class="mt-0.5 line-clamp-1 text-[11px] leading-snug text-(--text-muted) sm:line-clamp-2 sm:text-xs"
			title={container.image}
		>
			{container.image}
		</p>

		{#if container.ports.length > 0}
			<div class="mt-1.5 flex flex-wrap gap-1">
				{#each container.ports as port (portLabel(port))}
					<span
						class="rounded bg-(--surface-inset) px-1.5 py-px font-mono text-[10px] text-(--text-secondary)"
					>
						{portLabel(port)}
					</span>
				{/each}
			</div>
		{/if}
	</div>

	<p
		class="truncate border-t border-(--border-subtle) pt-1.5 text-[10px] text-(--text-muted) group-hover:text-(--text-secondary) sm:text-[11px]"
	>
		{link ?? container.status}
	</p>
</a>
