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
		if (container.state !== 'running') return 'bg-white/25';
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
	class="group flex flex-col gap-4 rounded-lg border border-white/6 bg-white/4 p-5 backdrop-blur-lg transition {link
		? 'hover:border-white/12 hover:bg-white/[0.07]'
		: 'cursor-default'}"
>
	<div class="flex items-start justify-between gap-3">
		{#if icon && !iconFailed}
			<img
				src={icon}
				alt=""
				class="h-12 w-12 rounded-lg object-contain"
				loading="lazy"
				onerror={() => (iconFailed = true)}
			/>
		{:else}
			<div
				class="flex h-12 w-12 items-center justify-center rounded-lg bg-white/5 text-blue-200/40"
				aria-hidden="true"
			>
				<CubeIcon size={26} weight="duotone" />
			</div>
		{/if}

		<span
			class="flex items-center gap-1.5 rounded-full bg-white/4 px-2 py-1 text-[11px] text-blue-200/40"
		>
			<span class="h-1.5 w-1.5 rounded-full {tone}" aria-hidden="true"></span>
			{container.state}
		</span>
	</div>

	<div class="min-w-0 flex-1">
		<h3 class="truncate text-sm font-semibold text-blue-100/90">{title}</h3>
		<p class="mt-1 line-clamp-2 text-sm leading-snug text-blue-200/35" title={container.image}>
			{container.image}
		</p>

		{#if container.ports.length > 0}
			<div class="mt-3 flex flex-wrap gap-1">
				{#each container.ports as port (portLabel(port))}
					<span class="rounded bg-white/5 px-1.5 py-0.5 font-mono text-[11px] text-blue-200/45">
						{portLabel(port)}
					</span>
				{/each}
			</div>
		{/if}
	</div>

	<p
		class="truncate border-t border-white/6 pt-3 text-xs text-blue-200/25 group-hover:text-blue-200/45"
	>
		{link ?? container.status}
	</p>
</a>
