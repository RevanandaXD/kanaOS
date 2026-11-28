<script lang="ts">
	import { browser } from '$app/env';
	import { weatherLocationStore } from '#lib/weatherLocation.svelte.ts';

	type Current = {
		temperature_2m: number;
		relative_humidity_2m: number;
		wind_speed_10m: number;
		weather_code: number;
	};

	type WeatherCache = {
		latitude: number;
		longitude: number;
		timestamp: number;
		current: Current;
	};

	const CACHE_KEY = 'kanaos-weather-cache';
	const CACHE_TTL = 15 * 60 * 1000;

	let current = $state<Current | null>(null);
	let error = $state('');

	function readCache(latitude: number, longitude: number): WeatherCache | null {
		if (!browser) return null;
		try {
			const raw = sessionStorage.getItem(CACHE_KEY);
			if (!raw) return null;
			const parsed = JSON.parse(raw) as Partial<WeatherCache>;
			if (
				parsed.latitude === latitude &&
				parsed.longitude === longitude &&
				typeof parsed.timestamp === 'number' &&
				parsed.current &&
				typeof parsed.current.temperature_2m === 'number'
			) {
				return parsed as WeatherCache;
			}
		} catch {
			// ignore invalid cache
		}
		return null;
	}

	function writeCache(latitude: number, longitude: number, data: Current) {
		if (!browser) return;
		try {
			const payload: WeatherCache = {
				latitude,
				longitude,
				timestamp: Date.now(),
				current: data
			};
			sessionStorage.setItem(CACHE_KEY, JSON.stringify(payload));
		} catch {
			// ignore storage quota errors
		}
	}

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

	$effect(() => {
		const { latitude, longitude } = weatherLocationStore.current;
		const controller = new AbortController();

		const cached = readCache(latitude, longitude);
		if (cached) {
			current = cached.current;
			error = '';
		} else {
			current = null;
		}

		const load = async () => {
			try {
				const url =
					`https://api.open-meteo.com/v1/forecast?latitude=${latitude}&longitude=${longitude}` +
					`&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&timezone=auto`;
				const res = await fetch(url, { signal: controller.signal });
				if (!res.ok) throw new Error(String(res.status));
				const data = (await res.json()).current as Current;
				current = data;
				error = '';
				writeCache(latitude, longitude, data);
			} catch (err) {
				if (err instanceof DOMException && err.name === 'AbortError') return;
				if (!current) error = 'Weather unavailable';
			}
		};

		let timeoutId: ReturnType<typeof setTimeout> | undefined;
		let intervalId: ReturnType<typeof setInterval> | undefined;

		const age = cached ? Date.now() - cached.timestamp : CACHE_TTL;
		if (age >= CACHE_TTL) {
			void load();
			intervalId = setInterval(load, CACHE_TTL);
		} else {
			timeoutId = setTimeout(() => {
				void load();
				intervalId = setInterval(load, CACHE_TTL);
			}, CACHE_TTL - age);
		}

		return () => {
			controller.abort();
			if (timeoutId) clearTimeout(timeoutId);
			if (intervalId) clearInterval(intervalId);
		};
	});
</script>

<div
	class="rounded-xl border border-(--border-subtle) bg-(--surface-card) p-4 shadow-sm backdrop-blur-xl sm:p-5"
>
	<p class="truncate text-xs font-medium tracking-wide text-(--text-muted) uppercase">
		{weatherLocationStore.current.name}
	</p>

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
