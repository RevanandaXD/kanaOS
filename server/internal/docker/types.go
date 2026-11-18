package docker

type Port struct {
	Private uint16 `json:"private"`
	Public  uint16 `json:"public,omitempty"`
	Type    string `json:"type"`
}

type Container struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`
	Health  string `json:"health"`
	Status  string `json:"status"`
	Project string `json:"project,omitempty"`
	Service string `json:"service,omitempty"`
	Ports   []Port `json:"ports"`
	Created int64  `json:"created"`
}
