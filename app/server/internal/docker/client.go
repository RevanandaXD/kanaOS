package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

const (
	labelProject = "com.docker.compose.project"
	labelService = "com.docker.compose.service"
)

type Client struct {
	api *client.Client
}

// Koneksi Daemon Docker
func New() (*Client, error) {
	api, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, fmt.Errorf("membuat docker client: %w", err)
	}

	return &Client{api: api}, nil
}

func (c *Client) Close() error {
	return c.api.Close()
}

// Ping
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.api.Ping(ctx, client.PingOptions{})
	return err
}

func (c *Client) List(ctx context.Context) ([]Container, error) {
	res, err := c.api.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("list container: %w", err)
	}

	containers := make([]Container, 0, len(res.Items))
	for _, item := range res.Items {
		containers = append(containers, toContainer(item))
	}

	return containers, nil
}
