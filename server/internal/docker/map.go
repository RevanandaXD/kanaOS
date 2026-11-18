package docker

import (
	"fmt"
	"strings"

	"github.com/moby/moby/api/types/container"
)

func toContainer(item container.Summary) Container {
	name := ""
	if len(item.Names) > 0 {
		name = strings.TrimPrefix(item.Names[0], "/")
	}

	return Container{
		ID:      item.ID,
		Name:    name,
		Image:   item.Image,
		State:   string(item.State),
		Health:  healthOf(item),
		Status:  item.Status,
		Project: item.Labels[labelProject],
		Service: item.Labels[labelService],
		Ports:   portsOf(item.Ports),
		Created: item.Created,
	}
}

func healthOf(item container.Summary) string {
	if item.Health != nil {
		return string(item.Health.Status)
	}
	if h := healthFromStatus(item.Status); h != "" {
		return h
	}
	return string(container.NoHealthcheck)
}

func healthFromStatus(status string) string {
	switch {
	case strings.Contains(status, "(unhealthy)"):
		return string(container.Unhealthy)
	case strings.Contains(status, "(healthy)"):
		return string(container.Healthy)
	case strings.Contains(status, "(health: starting)"):
		return string(container.Starting)
	default:
		return ""
	}
}

func portsOf(in []container.PortSummary) []Port {
	ports := make([]Port, 0, len(in))
	seen := make(map[string]struct{}, len(in))

	for _, p := range in {
		key := fmt.Sprintf("%d/%s/%d", p.PrivatePort, p.Type, p.PublicPort)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		ports = append(ports, Port{
			Private: p.PrivatePort,
			Public:  p.PublicPort,
			Type:    p.Type,
		})
	}

	return ports
}
