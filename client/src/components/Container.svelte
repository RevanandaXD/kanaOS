<script lang="ts">
	import type { Service } from '../api/services';

	let { service }: { service: Service } = $props();
</script>

<a
	href={service.url}
	target="_blank"
	rel="noreferrer"
	class="group flex flex-col gap-4 rounded-lg border border-white/6 bg-white/4 p-5 backdrop-blur-lg transition hover:border-white/12 hover:bg-white/[0.07]"
>
	{#if service.icon}
		<img src={service.icon} alt="" class="h-12 w-12 rounded-lg object-contain" loading="lazy" />
	{:else}
		<div
			class="flex h-12 w-12 items-center justify-center rounded-lg bg-white/5 text-lg font-medium text-blue-200/60"
			aria-hidden="true"
		>
			{service.name.charAt(0).toUpperCase()}
		</div>
	{/if}

	<div class="min-w-0 flex-1">
		<h3 class="truncate text-sm font-semibold text-blue-100/90">{service.name}</h3>
		{#if service.description}
			<p class="mt-1 line-clamp-2 text-sm leading-snug text-blue-200/35">{service.description}</p>
		{/if}
	</div>

	<p class="truncate border-t border-white/6 pt-3 text-xs text-blue-200/25 group-hover:text-blue-200/45">
		{service.url}
	</p>
</a>
