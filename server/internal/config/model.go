package config

type Service struct {
	Name					string `yaml:"name" json:"name"`
	Description		string `yaml:"description" json:"description"`
	URL						string `yaml:"url" json:"url"`
	Icon					string `yaml:"icon" json:"icon"`
}

type Category struct {
	name					string 		`yaml:"name" json:"name"`
	Services			[]Service `yaml:"services" json:"services"`
}

type Config struct {
	Categories 		[]Category `yaml:"categories" json:"categories"`
}