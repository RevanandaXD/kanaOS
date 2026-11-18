import { browser } from '$app/env';

export type Theme = 'dark' | 'light';

const STORAGE_KEY = 'kanaos-theme';
const DEFAULT_THEME: Theme = 'dark';

function loadTheme(): Theme {
	if (!browser) return DEFAULT_THEME;
	const saved = localStorage.getItem(STORAGE_KEY);
	return saved === 'light' || saved === 'dark' ? saved : DEFAULT_THEME;
}

export function applyTheme(theme: Theme) {
	if (!browser) return;
	document.documentElement.dataset.theme = theme;
	document.documentElement.classList.toggle('dark', theme === 'dark');
	document.documentElement.classList.toggle('light', theme === 'light');
	document.documentElement.style.colorScheme = theme;
}

class ThemeStore {
	current = $state<Theme>(DEFAULT_THEME);

	constructor() {
		if (browser) {
			this.current = loadTheme();
			applyTheme(this.current);
		}
	}

	setTheme(theme: Theme) {
		this.current = theme;
		if (browser) {
			localStorage.setItem(STORAGE_KEY, theme);
			applyTheme(theme);
		}
	}

	toggle() {
		this.setTheme(this.current === 'dark' ? 'light' : 'dark');
	}
}

export const themeStore = new ThemeStore();
