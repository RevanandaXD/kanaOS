<div align="center">
  <img src="./public/Logo.png" alt="kanaOS logo" width="180" />
  <h1>kanaOS</h1>
  <p><strong>Simple minimalism dashboard for self-hosted services.</strong></p>

  <p>
    <a href="https://github.com/RevanandaXD/kanaOS/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/RevanandaXD/kanaOS/ci.yml?branch=main&label=CI&logo=github&style=for-the-badge" alt="CI" /></a>
    <a href="https://github.com/RevanandaXD/kanaOS/actions/workflows/docker.yml"><img src="https://img.shields.io/github/actions/workflow/status/RevanandaXD/kanaOS/docker.yml?branch=main&label=Docker&logo=docker&style=for-the-badge" alt="Docker" /></a>
    <a href="https://github.com/RevanandaXD/kanaOS/stargazers"><img src="https://img.shields.io/github/stars/RevanandaXD/kanaOS?logo=github&style=for-the-badge" alt="Stars" /></a>
    <a href="https://github.com/RevanandaXD/kanaOS/commits/main"><img src="https://img.shields.io/github/commit-activity/t/RevanandaXD/kanaOS?label=commits&style=for-the-badge" alt="Total commits" /></a>
    <a href="./LICENSE"><img src="https://img.shields.io/github/license/RevanandaXD/kanaOS?style=for-the-badge" alt="License" /></a>
  </p>

  <p>
    <img src="./public/Light-mode.png" alt="kanaOS light mode" width="48%" />
    <img src="./public/Dark-mode.png" alt="kanaOS dark mode" width="48%" />
  </p>
</div>

## Daftar isi

- [Tentang](#tentang)
- [Fitur](#fitur)
- [Statistik repository](#statistik-repository)
- [Teknologi](#teknologi)
- [Struktur proyek](#struktur-proyek)
- [Menjalankan secara lokal](#menjalankan-secara-lokal)
- [Menjalankan dengan Docker](#menjalankan-dengan-docker)
- [Konfigurasi](#konfigurasi)
- [Endpoint API](#endpoint-api)
- [Berkontribusi](#berkontribusi)
- [Lisensi](#lisensi)

## Tentang

kanaOS adalah dashboard minimalis untuk memantau server dan container Docker. Antarmukanya dibangun dengan SvelteKit, sedangkan backend-nya adalah API Go ringan yang mengambil data sistem dan Docker, lalu melakukan proxy ke frontend dalam satu proses container.

Cocok untuk home server / self-host yang ingin melihat kondisi mesin tanpa antarmuka yang berat.

## Fitur

- **Monitoring Docker real-time** — daftar container, status `running` / `stopped` / `other`, health check, dan ringkasan yang diperbarui tiap 5 detik.
- **Informasi sistem** — OS, hostname, kernel, arsitektur, uptime, penggunaan CPU (gauge), RAM, dan disk dari endpoint `/api/system`.
- **Widget dashboard** — jam, kalender, dan cuaca (Open-Meteo, koordinat dapat diubah di `Weather.svelte`).
- **Tema terang & gelap** beserta pilihan wallpaper, dengan pengaturan tersimpan di browser.
- **Ringan & satu container** — frontend dan API dijalankan bersama oleh `entrypoint.sh`, diekspos lewat satu port.
- **CI otomatis** — `gofmt`, `go vet`, `go build`, `go test`, `pnpm check`, dan `pnpm build` berjalan di GitHub Actions.
- **Statistik repository otomatis** — kartu SVG di bawah ini diperbarui setiap ada push dan tiap hari.

## Statistik repository

<p align="center">
  <img src="./public/stats.svg" alt="Statistik repository kanaOS: commit, contributor, star, fork, branch, dan tag" width="620" />
</p>

Kartu di atas menampilkan jumlah **commit**, **contributor**, **star**, **fork**, **branch**, **tag**, jumlah baris kode, serta tanggal commit pertama dan terakhir. File `public/stats.svg` dihasilkan oleh [`scripts/generate-stats.sh`](./scripts/generate-stats.sh) dan di-commit otomatis oleh [`.github/workflows/stats.yml`](./.github/workflows/stats.yml):

- **Trigger:** push ke `main` dan jadwal harian pukul 03.23 UTC (`workflow_dispatch` juga tersedia).
- **Sumber data:** `git rev-list/shortlog/branch/tag` untuk metrik lokal, GitHub REST API via `gh` untuk star, fork, watcher, dan open issue.
- **Anti-loop:** commit workflow memakai `[skip ci]`, actor bot diabaikan, dan perubahan pada `public/stats.svg` tidak memicu workflow lagi.

Badge di bagian atas README juga "live" karena langsung dibaca dari GitHub API (shield.io), bukan angka statis.

Untuk membuat ulang kartu secara manual:

```bash
bash scripts/generate-stats.sh
# atau tentukan lokasi output
bash scripts/generate-stats.sh public/stats.svg

# opsional: timpa angka GitHub tanpa gh/token
STARS=12 FORKS=3 WATCHERS=1 OPEN_ISSUES=0 bash scripts/generate-stats.sh
```

> Metrik commit dihitung dari seluruh riwayat git (`git rev-list --count HEAD`), jadi checkout harus menggunakan `fetch-depth: 0` seperti di workflow.

## Teknologi

| Bagian | Stack |
| --- | --- |
| Frontend | SvelteKit 3, Svelte 5 (runes), TypeScript, Tailwind CSS 4, Vite 8, Phosphor Icons |
| Backend | Go 1.27, `net/http`, Docker client (`github.com/moby/moby/client`), YAML |
| Build | pnpm 12, Node 26, Docker multi-stage, GHCR |
| Otomatisasi | GitHub Actions (CI, Docker, statistik repository) |

## Struktur proyek

```text
kanaOS/
├── app/
│   ├── client/              # SvelteKit frontend
│   │   └── src/
│   │       ├── api/         # helper fetch (/api/containers, /api/system)
│   │       ├── components/  # Clock, Calendar, Weather, Container, Gauge, Navbar, Dropdown
│   │       └── routes/      # halaman dashboard + layout & tema
│   └── server/              # API Go
│       └── internal/
│           ├── api/         # router & handler
│           ├── config/      # loader service.yml
│           ├── docker/      # client, cache, mapping, tipe data
│           └── system/      # pengumpul info sistem per platform
├── config/                  # service.yml (opsional)
├── public/                  # Logo.png, screenshot tema, stats.svg (generated)
├── scripts/                 # generate-stats.sh
└── .github/workflows/       # ci.yml, docker.yml, stats.yml
```

## Menjalankan secara lokal

### Prasyarat

- Node.js 26 dan pnpm 12
- Go 1.27
- DockerEngine / Docker socket yang bisa diakses
- Git Bash atau shell POSIX untuk menjalankan script statistik

```bash
git clone https://github.com/RevanandaXD/kanaOS.git
cd kanaOS
```

Jalankan frontend dan backend di dua terminal terpisah:

```bash
# Terminal 1 — frontend (http://localhost:3000/app  atau port yang ditampilkan Vite)
cd app/client
pnpm install
pnpm dev
```

```bash
# Terminal 2 — API (http://localhost:8080)
cd app/server
go run ./cmd/kanaos
```

Frontend memanggil endpoint relatif (`/api/...`) dan API melakukan reverse proxy ke `http://127.0.0.1:3000`, jadi port default sudah saling cocok.

Pemeriksaan sebelum commit:

```bash
cd app/client && pnpm check && pnpm lint && pnpm build
cd ../server && gofmt -l . && go vet ./... && go test ./... && go build ./...
```

## Menjalankan dengan Docker

Cara tercepat memakai image yang sudah dibangun di GHCR:

```bash
docker compose up -d
```

Atau build sendiri dari sumber:

```bash
docker compose up -d --build
```

Lalu buka <http://localhost:3000>. Compose sudah melakukan mount `/var/run/docker.sock` (read-only) agar API bisa membaca container, dan mount `./config` sebagai konfigurasi.

## Konfigurasi

Semua env bersifat opsional:

| Variabel | Default | Keterangan |
| --- | --- | --- |
| `PORT` | `8080` | Port API / container saat produksi |
| `WEB_UPSTREAM` | `http://127.0.0.1:3000` | Tujuan reverse proxy untuk SvelteKit |
| `CONFIG_PATH` | `../config/service.yml` | Lokasi file konfigurasi service (YAML) |
| `DOCKER_CACHE_TTL` | `3s` | Durasi cache hasil `docker ps` |

Contoh `config/service.yml`:

```yaml
categories:
  - name: Media
    services:
      - name: Jellyfin
        description: Media server
        url: http://jellyfin.local
        icon: jellyfin
```

File tidak harus ada; API akan tetap berjalan dengan daftar kategori kosong.

## Endpoint API

| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/health` | Health check (`{"status":"ok","name":"KanaOS API"}`) |
| `GET` | `/api/containers` | Daftar container Docker beserta status & health (cache TTL) |
| `GET` | `/api/system` | Info OS, hostname, kernel, arsitektur, uptime, CPU, RAM, disk |
| `GET` | `/api/services` | Isi konfigurasi service dari `service.yml` |
| `GET` | `/*` | Reverse proxy ke frontend SvelteKit |

Contoh:

```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/containers
curl http://localhost:8080/api/system
```

## Berkontribusi

1. Fork repository lalu buat branch dari `main`, misalnya `feat/nama-fitur`.
2. Ikuti gaya commit yang sudah dipakai: `feat(...)`, `fix(...)`, `refactor(...)`, `chore(...)`.
3. Pastikan CI hijau (`pnpm check`, `pnpm build`, `gofmt`, `go vet`, `go test`).
4. Buat pull request dengan deskripsi perubahan yang jelas.

Statistik repository tidak perlu diedit manual — workflow akan memperbaruinya sendiri.

## Lisensi

Proyek ini dirilis di bawah [GNU General Public License v3.0](./LICENSE) © Revananda.
