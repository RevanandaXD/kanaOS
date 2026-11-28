import { browser } from '$app/env';

export type WeatherLocation = {
	name: string;
	country?: string;
	latitude: number;
	longitude: number;
};

const STORAGE_KEY = 'kanaos-weather-location';

export const DEFAULT_WEATHER_LOCATION: WeatherLocation = {
	name: 'Jakarta',
	country: 'ID',
	latitude: -6.2,
	longitude: 106.8
};

function isValidLocation(value: unknown): value is WeatherLocation {
	if (!value || typeof value !== 'object') return false;
	const loc = value as Record<string, unknown>;
	return (
		typeof loc.name === 'string' &&
		loc.name.trim().length > 0 &&
		typeof loc.latitude === 'number' &&
		Number.isFinite(loc.latitude) &&
		typeof loc.longitude === 'number' &&
		Number.isFinite(loc.longitude) &&
		(loc.country === undefined || typeof loc.country === 'string')
	);
}

function loadLocation(): WeatherLocation {
	if (!browser) return DEFAULT_WEATHER_LOCATION;
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return DEFAULT_WEATHER_LOCATION;
		const parsed: unknown = JSON.parse(raw);
		return isValidLocation(parsed) ? parsed : DEFAULT_WEATHER_LOCATION;
	} catch {
		return DEFAULT_WEATHER_LOCATION;
	}
}

class WeatherLocationStore {
	current = $state<WeatherLocation>(DEFAULT_WEATHER_LOCATION);

	constructor() {
		if (browser) {
			this.current = loadLocation();
		}
	}

	setLocation(location: WeatherLocation) {
		this.current = location;
		if (browser) {
			localStorage.setItem(STORAGE_KEY, JSON.stringify(location));
		}
	}

	reset() {
		this.setLocation(DEFAULT_WEATHER_LOCATION);
	}
}

export const weatherLocationStore = new WeatherLocationStore();

