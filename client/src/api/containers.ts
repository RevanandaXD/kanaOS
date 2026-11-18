export interface Port {
	private: number;
	public?: number;
	type: string;
}

export interface Container {
	id: string;
	name: string;
	image: string;
	state: string;
	health: string;
	status: string;
	project?: string;
	service?: string;
	ports: Port[];
	created: number;
}

export interface ContainerResponse {
	containers: Container[];
}

export async function getContainers(): Promise<ContainerResponse> {
	const response = await fetch('/api/containers');

	if (!response.ok) {
		const body = (await response.json().catch(() => null)) as { error?: string } | null;
		throw new Error(body?.error ?? `Failed to fetch containers (${response.status})`);
	}

	return response.json();
}
