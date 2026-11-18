import { browser } from '$app/env';

export type WallpaperTransition = 'fade' | 'tiles' | 'wave';

export type WallpaperOption = {
	id: string;
	name: string;
	src: string;
};

export const wallpapers: WallpaperOption[] = [
	{ id: 'default1', name: 'Default 1', src: '/wallpaper-default1.jpg' },
	{ id: 'default2', name: 'Default 2', src: '/wallpaper-default2.jpg' },
	{ id: 'default3', name: 'Default 3', src: '/wallpaper-default3.jpg' },
	{ id: 'default4', name: 'Default 4', src: '/wallpaper-default4.jpg' },
	{ id: 'default5', name: 'Default 5', src: '/wallpaper-default5.jpg' },
	{ id: 'default6', name: 'Default 6', src: '/wallpaper-default6.jpg' },
	{ id: 'default7', name: 'Default 7', src: '/wallpaper-default7.jpg' },
	{ id: 'default8', name: 'Default 8', src: '/wallpaper-default8.jpg' },
];

const TRANSITIONS: WallpaperTransition[] = ['fade', 'tiles', 'wave'];
const STORAGE_KEY = 'kanaos-wallpaper';
const DEFAULT_WALLPAPER = wallpapers[0].src;

function applyWallpaper(src: string) {
	if (!browser) return;
	document.documentElement.style.setProperty('--kanaos-wallpaper', `url('${src}')`);
}

function loadWallpaper(): string {
	if (!browser) return DEFAULT_WALLPAPER;
	const saved = localStorage.getItem(STORAGE_KEY);
	return wallpapers.some((w) => w.src === saved) ? (saved as string) : DEFAULT_WALLPAPER;
}

function randomTransition(): WallpaperTransition {
	return TRANSITIONS[Math.floor(Math.random() * TRANSITIONS.length)];
}

class WallpaperStore {
	current = $state(DEFAULT_WALLPAPER);
	previous = $state<string | null>(null);
	transition = $state<WallpaperTransition>('fade');
	changing = $state(false);
	#timer: ReturnType<typeof setTimeout> | null = null;

	constructor() {
		if (browser) {
			this.current = loadWallpaper();
			applyWallpaper(this.current);
		}
	}

	setWallpaper(src: string) {
		if (!browser) return;
		if (src === this.current) return;
		if (this.#timer) clearTimeout(this.#timer);
		this.transition = randomTransition();
		this.previous = this.current;
		this.current = src;
		this.changing = true;
		localStorage.setItem(STORAGE_KEY, src);
		applyWallpaper(src);
		const duration = this.transition === 'tiles' ? 900 : 700;
		this.#timer = setTimeout(() => {
			this.previous = null;
			this.changing = false;
			this.#timer = null;
		}, duration);
	}
}

export const wallpaperStore = new WallpaperStore();
