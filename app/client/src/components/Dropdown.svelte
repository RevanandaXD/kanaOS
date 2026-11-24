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
	import Gauge from './Gauge.svelte';
	import { wallpaperStore, wallpapers } from '#lib/wallpaper.svelte.ts';
	import { themeStore } from '#lib/theme.svelte.ts';
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
	class="flex h-[min(80vh,26rem)] w-full max-w-[min(92vw,45rem)] overflow-hidden rounded-xl border border-(--border-subtle) bg-(--surface-card) text-(--text-primary) shadow-2xl backdrop-blur-xl"
>
	<!-- Sidebar -->
	<aside
		class="flex w-13 shrink-0 flex-col items-center border-r border-(--border-subtle) bg-(--surface-inset) py-3 sm:w-14"
	>
		<nav class="flex flex-1 flex-col items-center gap-2">
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
	<section class="flex min-w-0 flex-1 flex-col">
		<!-- Wallpaper Section -->
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
							class="mt-3 inline-flex items-center gap-1.5 rounded-sm border border-rose-400/30 bg-white px-3 py-1.5 text-xs font-medium text-red-500 transition cursor-pointer"
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
					<div class="grid grid-cols-3 gap-2.5">
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
