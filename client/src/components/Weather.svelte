<script lang="ts">
	import { onMount } from 'svelte';

	const LAT = -6.2;
	const LON = 106.8;
	const PLACE = 'Lokasi saya';

	type Current = {
		temperature_2m: number;
		relative_humidity_2m: number;
		wind_speed_10m: number;
		weather_code: number;
	};

	let current = $state<Current | null>(null);
	let error = $state('');

	const labels: Record<number, string> = {
		0: 'Cerah', 1: 'Cerah berawan', 2: 'Berawan', 3: 'Mendung',
		45: 'Berkabut', 48: 'Berkabut',
		51: 'Gerimis', 53: 'Gerimis', 55: 'Gerimis',
		61: 'Hujan ringan', 63: 'Hujan', 65: 'Hujan lebat',
		80: 'Hujan lokal', 81: 'Hujan lokal', 82: 'Hujan deras',
		95: 'Badai petir', 96: 'Badai petir', 99: 'Badai petir'
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
				error = 'Cuaca tidak tersedia';
			}
		};
		load();
		const id = setInterval(load, 15 * 60 * 1000);
		return () => clearInterval(id);
	});
</script>

<div class="rounded-lg border border-white/[0.06] bg-white/[0.04] p-5 backdrop-blur-lg">
	<p class="text-xs font-medium tracking-wide text-blue-200/40 uppercase">{PLACE}</p>

	{#if error}
		<p class="mt-3 text-sm text-blue-200/35">{error}</p>
	{:else if !current}
		<p class="mt-3 text-sm text-blue-200/25">Memuat…</p>
	{:else}
		<div class="mt-2 flex items-end justify-between">
			<span class="text-4xl font-light text-blue-100/90">{Math.round(current.temperature_2m)}°</span>
			<span class="pb-1 text-sm text-blue-200/50">{labels[current.weather_code] ?? '—'}</span>
		</div>
		<div class="mt-4 grid grid-cols-2 gap-2 text-xs text-blue-200/40">
			<div class="rounded-md bg-white/[0.03] px-3 py-2">
				Kelembapan <span class="block text-sm font-medium text-blue-100/80">{current.relative_humidity_2m}%</span>
			</div>
			<div class="rounded-md bg-white/[0.03] px-3 py-2">
				Angin <span class="block text-sm font-medium text-blue-100/80">{current.wind_speed_10m} km/j</span>
			</div>
		</div>
	{/if}
</div>
