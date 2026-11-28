<script lang="ts">
	import GearIcon from 'phosphor-svelte/lib/GearSix';
	import ImageIcon from 'phosphor-svelte/lib/Image';
	import DesktopIcon from 'phosphor-svelte/lib/DesktopIcon';
	import PaletteIcon from 'phosphor-svelte/lib/Palette';
	import InfoIcon from 'phosphor-svelte/lib/Info';
	import SunIcon from 'phosphor-svelte/lib/Sun';
	import MoonIcon from 'phosphor-svelte/lib/Moon';
	import CpuIcon from 'phosphor-svelte/lib/Cpu';
	import MemoryIcon from 'phosphor-svelte/lib/Memory';
	import HardDriveIcon from 'phosphor-svelte/lib/HardDrive';
	import ClockIcon from 'phosphor-svelte/lib/Clock';
	import EyeIcon from 'phosphor-svelte/lib/Eye';
	import EyeSlashIcon from 'phosphor-svelte/lib/EyeSlash';
	import MapPinIcon from 'phosphor-svelte/lib/MapPin';
	import MagnifyingGlassIcon from 'phosphor-svelte/lib/MagnifyingGlass';
	import Gauge from './Gauge.svelte';
	import { wallpaperStore, wallpapers } from '#lib/wallpaper.svelte.ts';
	import { themeStore } from '#lib/theme.svelte.ts';
	import { privacyStore } from '#lib/hidePort.svelte.ts';
	import { weatherLocationStore, type WeatherLocation } from '#lib/weatherLocation.svelte.ts';
	import { getSystem, type SystemInfo } from '../api/system';
	import { version, author } from '../../package.json' with { type: 'json' };

	type MenuItem = {
		id: string;
		label: string;
		icon: typeof ImageIcon;
	};

	const menuItems: MenuItem[] = [
		{ id: 'systemInformation', label: 'System Information', icon: DesktopIcon },
		{ id: 'wallpaper', label: 'Wallpaper', icon: ImageIcon },
		{ id: 'appearance', label: 'Appearance', icon: PaletteIcon },
		{ id: 'general', label: 'General', icon: GearIcon },
		{ id: 'about', label: 'About', icon: InfoIcon }
	];

	let active = $state('systemInformation');

	// styling driven by CSS theme variables

	let systemInfo = $state<SystemInfo | null>(null);
	let systemError = $state('');
	let systemLoading = $state(false);

	type GeocodingResult = WeatherLocation & {
		id: number;
		admin1?: string;
	};

	let locationQuery = $state('');
	let locationResults = $state<GeocodingResult[]>([]);
	let locationSearching = $state(false);
	let locationError = $state('');

	async function searchLocation(event?: SubmitEvent) {
		event?.preventDefault();
		const q = locationQuery.trim();
		if (q.length < 2) {
			locationResults = [];
			locationError = '';
			return;
		}
		locationSearching = true;
		locationError = '';
		try {
			const url = `https://geocoding-api.open-meteo.com/v1/search?name=${encodeURIComponent(q)}&count=5&language=en&format=json`;
			const res = await fetch(url);
			if (!res.ok) throw new Error(String(res.status));
			const data = (await res.json()) as { results?: GeocodingResult[] };
			locationResults = data.results ?? [];
			if (locationResults.length === 0) {
				locationError = 'No matching cities found';
			}
		} catch {
			locationError = 'Failed to search location';
			locationResults = [];
		} finally {
			locationSearching = false;
		}
	}

	function selectLocation(loc: GeocodingResult) {
		weatherLocationStore.setLocation({
			name: loc.name,
			country: loc.country,
			latitude: loc.latitude,
			longitude: loc.longitude
		});
		locationQuery = '';
		locationResults = [];
		locationError = '';
	}

	// Real-time metric history for chart (last 16 points)
	let cpuHistory = $state<number[]>([
		12, 18, 15, 24, 20, 28, 22, 35, 25, 30, 28, 32, 26, 38, 30, 25
	]);
	let ramHistory = $state<number[]>([
		40, 42, 41, 43, 42, 44, 45, 43, 44, 46, 45, 47, 46, 48, 47, 48
	]);

	function pushMetric(val: number, arr: number[]): number[] {
		const next = [...arr.slice(1), Math.max(2, Math.min(100, Math.round(val)))];
		return next;
	}

	function makeSvgPath(points: number[], width = 200, height = 45): { line: string; area: string } {
		if (!points.length) return { line: '', area: '' };
		const step = width / (points.length - 1);
		const coords = points.map((p, i) => {
			const x = i * step;
			const y = height - (p / 100) * (height - 6) - 3;
			return [x, y];
		});

		const line = coords.reduce(
			(acc, [x, y], i) => (i === 0 ? `M ${x} ${y}` : `${acc} L ${x} ${y}`),
			''
		);
		const area = `${line} L ${width} ${height} L 0 ${height} Z`;
		return { line, area };
	}

	async function loadSystem() {
		if (systemLoading) return;
		systemLoading = true;
		systemError = '';
		try {
			systemInfo = await getSystem();
			if (systemInfo) {
				const cUsage = systemInfo.cpu_usage ?? 0;
				const rUsage = systemInfo.mem_usage ?? percent(systemInfo.mem_used, systemInfo.mem_total);
				cpuHistory = pushMetric(cUsage, cpuHistory);
				ramHistory = pushMetric(rUsage, ramHistory);
			}
		} catch (err) {
			systemError = err instanceof Error ? err.message : 'Failed to load system information';
		} finally {
			systemLoading = false;
		}
	}

	$effect(() => {
		if (active === 'systemInformation' && !systemInfo && !systemLoading && !systemError) {
			void loadSystem();
		}
	});

	// Auto-refresh metrics every 3s when dropdown & system tab active
	$effect(() => {
		if (active !== 'systemInformation' || !systemInfo) return;
		const timer = setInterval(() => {
			void loadSystem();
		}, 3000);
		return () => clearInterval(timer);
	});

	function formatBytes(bytes: number): string {
		if (!bytes) return '0 B';
		const gib = bytes / 1024 ** 3;
		if (gib >= 1) return `${gib.toFixed(gib >= 10 ? 1 : 2)} GiB`;
		const mib = bytes / 1024 ** 2;
		return `${mib.toFixed(0)} MiB`;
	}

	function formatUptime(seconds: number): string {
		const d = Math.floor(seconds / 86400);
		const h = Math.floor((seconds % 86400) / 3600);
		const m = Math.floor((seconds % 3600) / 60);
		const parts: string[] = [];
		if (d) parts.push(`${d}d`);
		if (h) parts.push(`${h}h`);
		if (m || parts.length === 0) parts.push(`${m}m`);
		return parts.join(' ');
	}

	function percent(used: number, total: number): number {
		if (!total) return 0;
		return Math.min(100, Math.round((used / total) * 100));
	}
</script>

<div
	class="flex max-h-[min(78dvh,32rem)] w-full flex-col overflow-hidden rounded-xl border border-(--border-subtle) bg-(--surface-card) text-(--text-primary) shadow-2xl backdrop-blur-xl sm:h-[min(80vh,26rem)] sm:flex-row"
>
	<aside
		class="flex shrink-0 items-center border-b border-(--border-subtle) bg-(--surface-inset) px-2 py-1.5 sm:w-14 sm:flex-col sm:border-r sm:border-b-0 sm:px-0 sm:py-3"
	>
		<nav class="flex flex-1 flex-row items-center justify-center gap-2 sm:flex-col sm:justify-start">
			{#each menuItems as item (item.id)}
				{@const Icon = item.icon}
				<button
					type="button"
					onclick={() => (active = item.id)}
					title={item.label}
					aria-label={item.label}
					class="flex h-9 w-9 items-center justify-center rounded-lg transition-colors {active ===
					item.id
						? 'bg-(--accent-soft) text-(--accent)'
						: 'text-(--text-muted) hover:bg-(--surface-hover) hover:text-(--text-primary)'}"
				>
					<Icon size={18} weight={active === item.id ? 'fill' : 'regular'} />
				</button>
			{/each}
		</nav>
	</aside>

	<!-- Content -->
	<section class="flex min-h-0 min-w-0 flex-1 flex-col">
		{#if active === 'wallpaper'}
			<header
				class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4"
			>
				<ImageIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">Wallpaper</h2>
					<p class="text-xs text-(--text-secondary)">
						Choose a default background from the static folder
					</p>
				</div>
			</header>

			<div class="flex-1 overflow-y-auto p-4 sm:p-5">
				<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-4">
					{#each wallpapers as wp (wp.id)}
						<button
							type="button"
							onclick={() => wallpaperStore.setWallpaper(wp.src)}
							class="group overflow-hidden rounded-lg border text-left transition {wallpaperStore.current ===
							wp.src
								? 'border-(--accent)/50 bg-(--accent-soft) ring-1 ring-(--accent)/30'
								: 'border-(--border-subtle) bg-(--surface-inset) hover:border-(--text-muted)/40 hover:bg-(--surface-hover)'}"
						>
							<div class="aspect-video w-full overflow-hidden">
								<img
									src={wp.src}
									alt={wp.name}
									class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
									loading="lazy"
								/>
							</div>
						</button>
					{/each}
				</div>
			</div>

			<!-- Appearance Section -->
		{:else if active === 'appearance'}
			<header
				class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4"
			>
				<PaletteIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">Appearance</h2>
					<p class="text-xs text-(--text-secondary)">Customize the interface appearance</p>
				</div>
			</header>
			<div class="flex-1 space-y-4 overflow-y-auto p-4 text-sm text-(--text-secondary) sm:p-5">
				<div
					class="flex items-center justify-between rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3"
				>
					<div>
						<p class="text-(--text-primary)">Theme</p>
						<p class="text-xs text-(--text-secondary)">
							{themeStore.current === 'dark' ? 'Dark mode active' : 'Light mode active'}
						</p>
					</div>
					<button
						type="button"
						role="switch"
						aria-checked={themeStore.current === 'light'}
						aria-label="Toggle light/dark theme"
						onclick={() => themeStore.toggle()}
						class="relative inline-flex h-7 w-13 items-center rounded-full border border-(--border-subtle) transition-colors {themeStore.current ===
						'light'
							? 'bg-(--accent)'
							: 'bg-(--surface-hover)'}"
					>
						<span
							class="inline-flex h-5 w-5 items-center justify-center rounded-full bg-white text-slate-700 shadow transition-transform {themeStore.current ===
							'light'
								? 'translate-x-7'
								: 'translate-x-1'}"
						>
							{#if themeStore.current === 'light'}
								<SunIcon size={12} weight="fill" />
							{:else}
								<MoonIcon size={12} weight="fill" />
							{/if}
						</span>
					</button>
				</div>
				<div
					class="flex items-center justify-between rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3"
				>
					<div>
						<p class="text-(--text-primary)">Accent</p>
						<p class="text-xs text-(--text-secondary)">Primary accent color</p>
					</div>
					<span class="h-5 w-5 rounded-full bg-(--accent) ring-2 ring-(--border-subtle)"></span>
				</div>
			</div>

			<!-- General Section -->
		{:else if active === 'general'}
			<header
				class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4"
			>
				<GearIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">General</h2>
					<p class="text-xs text-(--text-secondary)">General system settings</p>
				</div>
			</header>
			<div class="flex-1 space-y-3 overflow-y-auto p-4 text-sm sm:p-5">
				<div
					class="flex items-center justify-between rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3"
				>
					<div>
						<p class="text-(--text-primary)">Hide Port</p>
						<p class="mt-0.5 text-xs text-(--text-secondary)">
							{privacyStore.current === 'hide'
								? 'Port numbers are hidden on container cards'
								: 'Port numbers are visible on container cards'}
						</p>
					</div>
					<button
						type="button"
						role="switch"
						aria-checked={privacyStore.current === 'hide'}
						aria-label="Hide container ports"
						onclick={() => privacyStore.toggle()}
						class="relative inline-flex h-7 w-13 shrink-0 items-center rounded-full border border-(--border-subtle) transition-colors {privacyStore.current ===
						'hide'
							? 'bg-(--accent)'
							: 'bg-(--surface-hover)'}"
					>
						<span
							class="inline-flex h-5 w-5 items-center justify-center rounded-full bg-white text-slate-700 shadow transition-transform {privacyStore.current ===
							'hide'
								? 'translate-x-7'
								: 'translate-x-1'}"
						>
							{#if privacyStore.current === 'hide'}
								<EyeSlashIcon size={12} weight="fill" />
							{:else}
								<EyeIcon size={12} weight="fill" />
							{/if}
						</span>
					</button>
				</div>
				<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
					<div class="flex items-center justify-between gap-2">
						<div class="min-w-0">
							<p class="text-(--text-primary)">Weather Location</p>
							<p class="mt-0.5 flex items-center gap-1 truncate text-xs text-(--text-secondary)">
								<MapPinIcon size={12} class="shrink-0" />
								<span class="truncate">
									{weatherLocationStore.current.name}{weatherLocationStore.current.country
										? `, ${weatherLocationStore.current.country}`
										: ''}
								</span>
							</p>
						</div>
					</div>
					<form onsubmit={searchLocation} class="mt-3 flex gap-2">
						<input
							type="search"
							bind:value={locationQuery}
							placeholder="Search city (e.g. Tokyo, Bandung)…"
							aria-label="Search weather city"
							class="min-w-0 flex-1 rounded-md border border-(--border-subtle) bg-(--surface-card) px-3 py-1.5 text-xs text-(--text-primary) placeholder:text-(--text-muted) focus:border-(--accent) focus:outline-none"
						/>
						<button
							type="submit"
							disabled={locationSearching}
							class="inline-flex shrink-0 items-center gap-1.5 rounded-md bg-(--accent-soft) px-3 py-1.5 text-xs font-medium text-(--accent) transition hover:opacity-90 disabled:opacity-50"
						>
							<MagnifyingGlassIcon size={14} />
							<span>{locationSearching ? 'Searching…' : 'Search'}</span>
						</button>
					</form>
					{#if locationError}
						<p class="mt-2 text-xs text-(--text-muted)">{locationError}</p>
					{:else if locationResults.length > 0}
						<ul class="mt-2 divide-y divide-(--border-subtle) rounded-md border border-(--border-subtle) bg-(--surface-card)">
							{#each locationResults as loc (loc.id)}
								<li>
									<button
										type="button"
										onclick={() => selectLocation(loc)}
										class="flex w-full items-center justify-between gap-2 px-3 py-2 text-left text-xs transition hover:bg-(--surface-hover)"
									>
										<span class="truncate font-medium text-(--text-primary)">
											{loc.name}{loc.admin1 ? `, ${loc.admin1}` : ''}
										</span>
										{#if loc.country}
											<span class="shrink-0 text-[11px] text-(--text-muted)">{loc.country}</span>
										{/if}
									</button>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
				<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
					<p class="text-(--text-primary)">Language</p>
					<p class="mt-0.5 text-xs text-(--text-secondary)">English</p>
				</div>
				<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
					<p class="text-(--text-primary)">Time zone</p>
					<p class="mt-0.5 text-xs text-(--text-secondary)">Automatic (system)</p>
				</div>
			</div>

			<!-- System Information Section -->
		{:else if active === 'systemInformation'}
			<header
				class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4"
			>
				<DesktopIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">System Information</h2>
					<p class="text-xs text-(--text-secondary)">Host and system resource information</p>
				</div>
			</header>

			<div class="flex-1 space-y-4 overflow-y-auto p-4 sm:p-5">
				{#if systemLoading && !systemInfo}
					<div class="flex h-40 flex-col items-center justify-center gap-2 text-(--text-secondary)">
						<div
							class="h-5 w-5 animate-spin rounded-full border-2 border-(--border-subtle) border-t-(--accent)"
						></div>
						<p class="text-xs">Loading system information…</p>
					</div>
				{:else if systemError && !systemInfo}
					<div class="rounded-xl border border-red-400/20 bg-red-500 p-4 text-white">
						<div class="flex items-center gap-2 text-sm font-medium text-rose-300">
							<svg class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
								/>
							</svg>
							<span>Unable to load system info</span>
						</div>
						<p class="mt-1.5 text-xs text-rose-200/70">{systemError}</p>
						<button
							type="button"
							class="mt-3 inline-flex cursor-pointer items-center gap-1.5 rounded-sm border border-rose-400/30 bg-white px-3 py-1.5 text-xs font-medium text-red-500 transition"
							onclick={() => {
								systemInfo = null;
								systemError = '';
								void loadSystem();
							}}
						>
							Retry connection
						</button>
					</div>
				{:else if systemInfo}
					<!-- System Specs Grid -->
					<div class="grid grid-cols-2 gap-2 sm:grid-cols-4 sm:gap-2.5">
						<div
							class="rounded-md border border-(--border-subtle) bg-(--surface-inset) p-3 transition hover:border-(--text-muted)/30"
						>
							<span class="text-[10px] font-semibold tracking-wider text-(--text-muted) uppercase"
								>OS</span
							>
							<p
								class="mt-1 truncate text-xs font-medium text-(--text-primary)"
								title={systemInfo.os}
							>
								{systemInfo.os || '-'}
							</p>
						</div>
						<div
							class="rounded-md border border-(--border-subtle) bg-(--surface-inset) p-3 transition hover:border-(--text-muted)/30"
						>
							<span class="text-[10px] font-semibold tracking-wider text-(--text-muted) uppercase"
								>Host</span
							>
							<p
								class="mt-1 truncate text-xs font-medium text-(--text-primary)"
								title={systemInfo.hostname}
							>
								{systemInfo.hostname || '-'}
							</p>
						</div>
						<div
							class="rounded-md border border-(--border-subtle) bg-(--surface-inset) p-3 transition hover:border-(--text-muted)/30"
						>
							<span class="text-[10px] font-semibold tracking-wider text-(--text-muted) uppercase"
								>Kernel</span
							>
							<p
								class="mt-1 truncate text-xs font-medium text-(--text-primary)"
								title={systemInfo.kernel}
							>
								{systemInfo.kernel || '-'}
							</p>
						</div>
						<div
							class="rounded-md border border-(--border-subtle) bg-(--surface-inset) p-3 transition hover:border-(--text-muted)/30"
						>
							<span class="text-[10px] font-semibold tracking-wider text-(--text-muted) uppercase"
								>Arch</span
							>
							<p
								class="mt-1 truncate text-xs font-medium text-(--text-primary)"
								title={systemInfo.arch}
							>
								{systemInfo.arch || '-'}
							</p>
						</div>
					</div>

					{@const disk = systemInfo.disks?.[0]}
					<div class="grid grid-cols-3 gap-2 sm:gap-2.5">
						<Gauge
							Icon={CpuIcon}
							label="CPU"
							value={systemInfo.cpu_usage ?? 0}
							sub="{systemInfo.cpus} cores"
						/>
						<Gauge
							Icon={MemoryIcon}
							label="Memory"
							color="#818cf8"
							value={systemInfo.mem_usage ?? percent(systemInfo.mem_used, systemInfo.mem_total)}
							sub="{formatBytes(systemInfo.mem_used)} / {formatBytes(systemInfo.mem_total)}"
						/>
						<Gauge
							Icon={HardDriveIcon}
							label="Disk"
							color="#34d399"
							value={disk ? percent(disk.used, disk.total) : 0}
							sub={disk ? `${formatBytes(disk.used)} / ${formatBytes(disk.total)}` : 'N/A'}
						/>
					</div>

					<div
						class="flex items-center justify-between rounded-md border border-(--border-subtle) bg-(--surface-inset) px-4 py-3 text-xs"
					>
						<div class="flex items-center gap-2 text-(--text-muted)">
							<ClockIcon size={14} />
							<span>System Uptime</span>
						</div>
						<span class="font-mono text-[11px] font-medium text-(--text-primary)">
							{formatUptime(systemInfo.uptime)}
						</span>
					</div>
				{/if}
			</div>

			<!-- About Section -->
		{:else}
			<header
				class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4"
			>
				<InfoIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">About</h2>
					<p class="text-xs text-(--text-secondary)">Application information</p>
				</div>
			</header>
			<div class="flex-1 overflow-y-auto p-4 sm:p-5">
				<div
					class="flex flex-col items-center gap-5 rounded-lg border border-(--border-subtle) bg-(--surface-inset) p-4 sm:flex-row sm:gap-6"
				>
					<img src="./Logo.png" alt="Logo KanaOS" class="h-35 w-auto shrink-0 object-contain" />

					<section class="min-w-0 text-center sm:text-left">
						<h2 class="text-base font-semibold text-(--text-primary)">KanaOS</h2>
						<p class="mt-1 text-xs text-(--text-secondary)">
							A personal dashboard for managing services and system appearance.
						</p>

						<p class="mt-4 text-xs text-(--text-muted)">Version {version} | Made by {author}</p>
					</section>
				</div>
			</div>
		{/if}
	</section>
</div>
