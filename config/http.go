package config

type ConfigHTTP struct {
	Address       string `yaml:"address"`
	Port          string `yaml:"port"`
	AllowInsecure bool   `yaml:"allow_insecure"`
}
