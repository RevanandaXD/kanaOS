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

	async function loadSystem() {
		if (systemLoading) return;
		systemLoading = true;
		systemError = '';
		try {
			systemInfo = await getSystem();
		} catch (err) {
			systemError = err instanceof Error ? err.message : 'Failed to load system information';
		} finally {
			systemLoading = false;
		}
	}

	$effect(() => {
		if (active === 'systemInformation' && !systemInfo && !systemLoading) {
			void loadSystem();
		}
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
	class="flex h-[min(80vh,26rem)] w-full max-w-[min(92vw,45rem)] overflow-hidden rounded-2xl border border-(--border-subtle) bg-(--surface-card) text-(--text-primary) shadow-2xl backdrop-blur-xl"
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
					class="flex h-9 w-9 items-center justify-center rounded-lg transition-colors {active === item.id
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
			<header class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4">
				<ImageIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">Wallpaper</h2>
					<p class="text-xs text-(--text-secondary)">Choose a default background from the static folder</p>
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
			<header class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4">
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
			<header class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4">
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
			<header class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4">
				<DesktopIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-(--text-primary)">System Information</h2>
					<p class="text-xs text-(--text-secondary)">Host and system resource information</p>
				</div>
			</header>

			<div class="flex-1 space-y-4 overflow-y-auto p-4 text-sm sm:p-5">
				{#if systemLoading && !systemInfo}
					<p class="text-(--text-secondary)">Loading system information…</p>
				{:else if systemError && !systemInfo}
					<div class="rounded-lg border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-rose-200/80">
						<p class="font-medium">Unable to load system information</p>
						<p class="mt-1 text-xs text-rose-200/60">{systemError}</p>
						<button
							type="button"
							class="mt-2 text-xs text-(--text-secondary) underline underline-offset-2 hover:text-(--text-primary)"
							onclick={() => {
								systemInfo = null;
								void loadSystem();
							}}
						>
							Try again
						</button>
					</div>
				{:else if systemInfo}
					<div class="grid gap-3 sm:grid-cols-2">
						<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
							<p class="text-[11px] tracking-wide text-(--text-secondary) uppercase">OS</p>
							<p class="mt-1 font-medium text-(--text-primary)" title={systemInfo.os}>{systemInfo.os || '-'}</p>
						</div>
						<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
							<p class="text-[11px] tracking-wide text-(--text-secondary) uppercase">Host</p>
							<p class="mt-1 font-medium text-(--text-primary)">{systemInfo.hostname || '-'}</p>
						</div>
						<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
							<p class="text-[11px] tracking-wide text-(--text-secondary) uppercase">Kernel</p>
							<p class="mt-1 font-medium text-(--text-primary)">{systemInfo.kernel || '-'}</p>
						</div>
						<div class="rounded-lg border border-(--border-subtle) bg-(--surface-inset) px-4 py-3">
							<p class="text-[11px] tracking-wide text-(--text-secondary) uppercase">Architecture</p>
							<p class="mt-1 font-medium text-(--text-primary)">{systemInfo.arch || '-'}</p>
						</div>
					</div>

					<div class="space-y-3 rounded-lg border border-(--border-subtle) bg-(--surface-inset) p-4">
						<div class="flex items-start gap-3">
							<CpuIcon size={18} class="mt-0.5 shrink-0 text-(--text-secondary)" />
							<div class="min-w-0 flex-1">
								<div class="flex items-baseline justify-between gap-2">
									<p class="text-(--text-primary)">CPU</p>
									<p class="text-xs text-(--text-secondary)">{systemInfo.cpus} core</p>
								</div>
								<p class="mt-0.5 truncate text-xs text-(--text-secondary)" title={systemInfo.cpu_model}>
									{systemInfo.cpu_model || 'Unknown model'}
								</p>
							</div>
						</div>

						<div class="flex items-start gap-3">
							<MemoryIcon size={18} class="mt-0.5 shrink-0 text-(--text-secondary)" />
							<div class="min-w-0 flex-1">
								<div class="flex items-baseline justify-between gap-2">
									<p class="text-(--text-primary)">Memory</p>
									<p class="text-xs text-(--text-secondary)">
										{formatBytes(systemInfo.mem_used)} / {formatBytes(systemInfo.mem_total)}
										<span class="text-(--text-muted)">({percent(systemInfo.mem_used, systemInfo.mem_total)}%)</span>
									</p>
								</div>
								<div class="mt-2 h-1.5 overflow-hidden rounded-full bg-(--surface-hover)">
									<div
										class="h-full rounded-full bg-(--accent)"
										style="width: {percent(systemInfo.mem_used, systemInfo.mem_total)}%"
									></div>
								</div>
							</div>
						</div>

						{#if systemInfo.disks?.length}
							{@const disk = systemInfo.disks[0]}
							<div class="flex items-start gap-3">
								<HardDriveIcon size={18} class="mt-0.5 shrink-0 text-(--text-secondary)" />
								<div class="min-w-0 flex-1">
									<div class="flex items-baseline justify-between gap-2">
										<p class="text-(--text-primary)">Disk ({disk.path})</p>
										<p class="text-xs text-(--text-secondary)">
											{formatBytes(disk.used)} / {formatBytes(disk.total)}
											<span class="text-(--text-muted)">({percent(disk.used, disk.total)}%)</span>
										</p>
									</div>
									<div class="mt-2 h-1.5 overflow-hidden rounded-full bg-(--surface-hover)">
										<div
											class="h-full rounded-full bg-emerald-400/60"
											style="width: {percent(disk.used, disk.total)}%"
										></div>
									</div>
								</div>
							</div>
						{/if}

						<div class="flex items-start gap-3">
							<ClockIcon size={18} class="mt-0.5 shrink-0 text-(--text-secondary)" />
							<div class="min-w-0 flex-1">
								<div class="flex items-baseline justify-between gap-2">
									<p class="text-(--text-primary)">Uptime</p>
									<p class="text-xs text-(--text-secondary)">{formatUptime(systemInfo.uptime)}</p>
								</div>
							</div>
						</div>
					</div>
				{/if}
			</div>

			<!-- About Section -->
		{:else}
			<header class="flex items-center gap-3 border-b border-(--border-subtle) px-4 py-3 sm:px-5 sm:py-4">
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
