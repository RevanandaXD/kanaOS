<script lang="ts">
	import GearIcon from 'phosphor-svelte/lib/GearSix';
	import ImageIcon from 'phosphor-svelte/lib/Image';
	import PaletteIcon from 'phosphor-svelte/lib/Palette';
	import InfoIcon from 'phosphor-svelte/lib/Info';
	import SunIcon from 'phosphor-svelte/lib/Sun';
	import MoonIcon from 'phosphor-svelte/lib/Moon';
	import { wallpaperStore, wallpapers } from '#lib/wallpaper.svelte.ts';
	import { themeStore } from '#lib/theme.svelte.ts';
	import { version, author } from '../../package.json' with { type: 'json' };

	type MenuItem = {
		id: string;
		label: string;
		icon: typeof ImageIcon;
	};

	const menuItems: MenuItem[] = [
		{ id: 'wallpaper', label: 'Wallpaper', icon: ImageIcon },
		{ id: 'appearance', label: 'Appearance', icon: PaletteIcon },
		{ id: 'general', label: 'General', icon: GearIcon },
		{ id: 'about', label: 'About', icon: InfoIcon }
	];

	let active = $state('wallpaper');

	const isLight = $derived(themeStore.current === 'light');
</script>

<div
	class="flex h-105 w-180 max-w-[90vw] overflow-hidden rounded-xl border shadow-2xl backdrop-blur-xl {isLight
		? 'border-slate-200/70 bg-white/90 text-slate-800'
		: 'border-white/10 bg-[#0f1424]/95 text-blue-100'}"
>
	<!-- Sidebar -->
	<aside
		class="flex w-14 shrink-0 flex-col items-center border-r py-3 {isLight
			? 'border-slate-200/70 bg-slate-100/70'
			: 'border-white/10 bg-white/3'}"
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
						? isLight
							? 'bg-slate-900/10 text-slate-900'
							: 'bg-white/10 text-blue-50'
						: isLight
							? 'text-slate-500 hover:bg-slate-900/5 hover:text-slate-800'
							: 'text-blue-200/60 hover:bg-white/6 hover:text-blue-100/90'}"
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
			<header class="flex items-center gap-3 border-b border-white/10 px-5 py-4">
				<ImageIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-blue-100/90">Wallpaper</h2>
					<p class="text-xs text-blue-200/40">Pilih gambar latar default dari folder static</p>
				</div>
			</header>

			<div class="flex-1 overflow-y-auto p-5">
				<div class="grid grid-cols-1 gap-2 sm:grid-cols-4">
					{#each wallpapers as wp (wp.id)}
						<button
							type="button"
							onclick={() => wallpaperStore.setWallpaper(wp.src)}
							class="group overflow-hidden rounded-lg border text-left transition {wallpaperStore.current ===
							wp.src
								? 'border-blue-300/50 bg-white/8 ring-1 ring-blue-300/30'
								: 'border-white/10 bg-white/3 hover:border-white/25 hover:bg-white/6'}"
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
			<header class="flex items-center gap-3 border-b border-white/10 px-5 py-4">
				<PaletteIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-blue-100/90">Appearance</h2>
					<p class="text-xs text-blue-200/40">Sesuaikan tampilan antarmuka</p>
				</div>
			</header>
			<div class="flex-1 space-y-4 overflow-y-auto p-5 text-sm text-blue-200/60">
				<div
					class="flex items-center justify-between rounded-lg border border-white/10 bg-white/3 px-4 py-3"
				>
					<div>
						<p class="text-blue-100/85">Tema</p>
						<p class="text-xs text-blue-200/40">
							{themeStore.current === 'dark' ? 'Mode gelap aktif' : 'Mode terang aktif'}
						</p>
					</div>
					<button
						type="button"
						role="switch"
						aria-checked={themeStore.current === 'light'}
						aria-label="Toggle tema light/dark"
						onclick={() => themeStore.toggle()}
						class="relative inline-flex h-7 w-13 items-center rounded-full border border-white/10 transition-colors {themeStore.current === 'light'
							? 'bg-blue-400/70'
							: 'bg-white/10'}"
					>
						<span
							class="inline-flex h-5 w-5 items-center justify-center rounded-full bg-white text-slate-700 shadow transition-transform {themeStore.current === 'light'
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
					class="flex items-center justify-between rounded-lg border border-white/10 bg-white/3 px-4 py-3"
				>
					<div>
						<p class="text-blue-100/85">Aksen</p>
						<p class="text-xs text-blue-200/40">Warna aksen utama</p>
					</div>
					<span class="h-5 w-5 rounded-full bg-blue-400/80 ring-2 ring-white/10"></span>
				</div>
			</div>

			<!-- General Section -->
		{:else if active === 'general'}
			<header class="flex items-center gap-3 border-b border-white/10 px-5 py-4">
				<GearIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-blue-100/90">General</h2>
					<p class="text-xs text-blue-200/40">Pengaturan umum sistem</p>
				</div>
			</header>
			<div class="flex-1 space-y-3 overflow-y-auto p-5 text-sm">
				<div class="rounded-lg border border-white/10 bg-white/3 px-4 py-3">
					<p class="text-blue-100/85">Bahasa</p>
					<p class="mt-0.5 text-xs text-blue-200/40">Indonesia</p>
				</div>
				<div class="rounded-lg border border-white/10 bg-white/3 px-4 py-3">
					<p class="text-blue-100/85">Zona waktu</p>
					<p class="mt-0.5 text-xs text-blue-200/40">Otomatis (sistem)</p>
				</div>
			</div>

			<!-- About Section -->
		{:else}
			<header class="flex items-center gap-3 border-b border-white/10 px-5 py-4">
				<InfoIcon size={20} />
				<div>
					<h2 class="text-sm font-semibold text-blue-100/90">About</h2>
					<p class="text-xs text-blue-200/40">Informasi aplikasi</p>
				</div>
			</header>
			<div class="flex-1 overflow-y-auto p-5">
				<div
					class="flex flex-col items-center gap-5 rounded-lg border border-white/10 bg-white/3 p-4 sm:flex-row sm:gap-6"
				>
					<img src="./Logo.png" alt="Logo KanaOS" class="h-35 w-auto shrink-0 object-contain" />

					<section class="min-w-0 text-center sm:text-left">
						<h2 class="text-base font-semibold text-blue-100/90">KanaOS</h2>
						<p class="mt-1 text-xs text-blue-200/50">
							Dashboard pribadi untuk mengelola layanan dan tampilan sistem.
						</p>

						<p class="mt-4 text-xs text-blue-200/30">Version {version} | Made by {author}</p>
					</section>
				</div>
			</div>
		{/if}
	</section>
</div>
