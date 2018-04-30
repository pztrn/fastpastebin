package config

type ConfigDatabase struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
