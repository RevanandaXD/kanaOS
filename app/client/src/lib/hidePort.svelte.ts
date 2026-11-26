import { browser } from '$app/env';

export type Privacy = 'hide' | 'show';

const STORAGE_KEY = 'ports-privacy';
const DEFAULT_PRIVACY: Privacy = 'hide';

function loadPrivacy(): Privacy {
	if (!browser) return DEFAULT_PRIVACY;
	const saved = localStorage.getItem(STORAGE_KEY);
	return saved === 'hide' || saved === 'show' ? saved : DEFAULT_PRIVACY;
}

export function applyPrivacy(privacy: Privacy) {
	if (!browser) return;
	document.documentElement.dataset.privacy = privacy;
	document.documentElement.classList.toggle('show', privacy === 'show');
	document.documentElement.classList.toggle('hide', privacy === 'hide');
}

class PrivacyStore {
	current = $state<Privacy>(DEFAULT_PRIVACY);

	constructor() {
		if (browser) {
			this.current = loadPrivacy();
			applyPrivacy(this.current);
		}
	}

	setPrivacy(privacy: Privacy) {
		this.current = privacy;
		if (browser) {
			localStorage.setItem(STORAGE_KEY, privacy);
			applyPrivacy(privacy);
		}
	}

	toggle() {
		this.setPrivacy(this.current === 'hide' ? 'show' : 'hide');
	}
}

export const privacyStore = new PrivacyStore();
