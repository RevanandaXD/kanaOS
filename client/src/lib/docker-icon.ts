const CDN = 'https://cdn.jsdelivr.net/gh/selfhst/icons@main/svg';

const aliases: Record<string, string> = {
	'portainer-ce': 'portainer',
	'portainer-ee': 'portainer',
	postgres: 'postgresql',
	mongo: 'mongodb',
	'immich-server': 'immich',
	homeassistant: 'home-assistant',
	adguardhome: 'adguard-home',
	'mariadb-server': 'mariadb'
};

export function dockerIcon(image: string): string | null {
	const repo = repoName(image);
	if (!repo) return null;

	const base = stripSuffix(repo);
	const slug = aliases[repo] ?? aliases[base] ?? repo ?? base;

	return `${CDN}/${slug}.svg`;
}

// "ghcr.io/immich-app/immich-server:release" -> "immich-server"
function repoName(image: string): string {
	const withoutDigest = image.split('@')[0];
	const withoutTag = withoutDigest.replace(/:[^/]*$/, '');

	return (withoutTag.split('/').pop() ?? '').toLowerCase();
}

// "portainer-ce" -> "portainer", "immich-server" -> "immich"
function stripSuffix(name: string): string {
	return name.replace(/-(server|ce|ee|app|web|api|ui)$/, '');
}
