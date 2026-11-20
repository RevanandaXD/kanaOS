<script lang="ts">
	import { onMount } from 'svelte';

	const LAT = -6.2;
	const LON = 106.8;
	const PLACE = 'My location';

	type Current = {
		temperature_2m: number;
		relative_humidity_2m: number;
		wind_speed_10m: number;
		weather_code: number;
	};

	let current = $state<Current | null>(null);
	let error = $state('');

	const labels: Record<number, string> = {
		0: 'Clear',
		1: 'Mostly clear',
		2: 'Partly cloudy',
		3: 'Overcast',
		45: 'Fog',
		48: 'Rime fog',
		51: 'Light drizzle',
		53: 'Drizzle',
		55: 'Heavy drizzle',
		61: 'Light rain',
		63: 'Rain',
		65: 'Heavy rain',
		80: 'Rain showers',
		81: 'Rain showers',
		82: 'Heavy showers',
		95: 'Thunderstorm',
		96: 'Thunderstorm',
		99: 'Thunderstorm'
	};

	onMount(() => {
		const load = async () => {
			try {
				const url =
					`https://api.open-meteo.com/v1/forecast?latitude=${LAT}&longitude=${LON}` +
					`&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&timezone=auto`;
				const res = await fetch(url);
				if (!res.ok) throw new Error(String(res.status));
				current = (await res.json()).current;
				error = '';
			} catch {
				error = 'Weather unavailable';
			}
		};
		load();
		const id = setInterval(load, 15 * 60 * 1000);
		return () => clearInterval(id);
	});
</script>

<div
	class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 shadow-sm backdrop-blur-xl sm:p-5"
>
	<p class="text-xs font-medium tracking-wide text-(--text-muted) uppercase">{PLACE}</p>

	{#if error}
		<p class="mt-3 text-sm text-(--text-muted)">{error}</p>
	{:else if !current}
		<p class="mt-3 text-sm text-(--text-muted)">Loading…</p>
	{:else}
		<div class="mt-2 flex items-end justify-between">
			<span class="text-3xl font-light text-(--text-primary) sm:text-4xl"
				>{Math.round(current.temperature_2m)}°</span
			>
			<span class="pb-1 text-sm text-(--text-secondary)">{labels[current.weather_code] ?? '—'}</span
			>
		</div>
		<div class="mt-4 grid grid-cols-2 gap-2 text-xs text-(--text-secondary)">
			<div class="rounded-md bg-(--surface-inset) px-3 py-2">
				Humidity <span class="mt-0.5 block text-sm font-medium text-(--text-primary)"
					>{current.relative_humidity_2m}%</span
				>
			</div>
			<div class="rounded-md bg-(--surface-inset) px-3 py-2">
				Wind <span class="mt-0.5 block text-sm font-medium text-(--text-primary)"
					>{current.wind_speed_10m} km/h</span
				>
			</div>
		</div>
	{/if}
</div>
