package docker

import (
	"context"
	"sync"
	"time"
)

// Cache menyimpan hasil list container selama TTL tertentu.
// Tujuannya bukan "biar cepat", tapi supaya jumlah panggilan ke Docker
// tidak bertambah sebanding dengan jumlah tab browser yang terbuka.
type Cache struct {
	client *Client
	ttl    time.Duration

	// fetchMu memastikan hanya satu goroutine yang memanggil Docker
	// dalam satu waktu. Lock ini terpisah dari mu supaya panggilan
	// jaringan tidak menahan lock data.
	fetchMu sync.Mutex

	mu         sync.RWMutex
	containers []Container
	fetchedAt  time.Time
	err        error
}

func NewCache(client *Client, ttl time.Duration) *Cache {
	return &Cache{client: client, ttl: ttl}
}

// Containers mengembalikan data dari cache bila masih segar.
// Bila sudah kedaluwarsa, satu goroutine akan refresh dan yang lain menunggu.
func (c *Cache) Containers(ctx context.Context) ([]Container, error) {
	if containers, err, ok := c.get(); ok {
		return containers, err
	}

	c.fetchMu.Lock()
	defer c.fetchMu.Unlock()

	// Cek ulang setelah dapat lock: bisa jadi goroutine lain baru saja
	// selesai refresh saat kita menunggu. Tanpa ini, N request yang
	// datang bersamaan akan tetap memicu N panggilan Docker.
	if containers, err, ok := c.get(); ok {
		return containers, err
	}

	containers, err := c.client.List(ctx)

	c.mu.Lock()
	c.containers = containers
	c.fetchedAt = time.Now()
	c.err = err
	c.mu.Unlock()

	return containers, err
}

func (c *Cache) get() ([]Container, error, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// fetchedAt masih zero value saat cache belum pernah terisi,
	// sehingga time.Since-nya sangat besar dan otomatis dianggap basi.
	if time.Since(c.fetchedAt) < c.ttl {
		return c.containers, c.err, true
	}

	return nil, nil, false
}
