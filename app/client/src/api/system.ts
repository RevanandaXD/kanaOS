export interface DiskInfo {
	path: string;
	total: number;
	used: number;
}

export interface SystemInfo {
	os: string;
	hostname: string;
	kernel: string;
	arch: string;
	uptime: number;
	cpus: number;
	cpu_model: string;
	mem_total: number;
	mem_used: number;
	disks: DiskInfo[];
}

export async function getSystem(): Promise<SystemInfo> {
	const response = await fetch('/api/system');

	if (!response.ok) {
		throw new Error(`Failed to fetch system info (${response.status})`);
	}

	return response.json();
}
