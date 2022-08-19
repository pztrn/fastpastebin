package application

// ConfigDatabase describes database configuration.
type ConfigDatabase struct {
	Type     string `yaml:"type"`
	Path     string `yaml:"path"`
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ConfigHTTP describes HTTP server configuration.
type ConfigHTTP struct {
	Address              string `yaml:"address"`
	Port                 string `yaml:"port"`
	MaxBodySizeMegabytes string `yaml:"max_body_size_megabytes"`
	AllowInsecure        bool   `yaml:"allow_insecure"`
}

// ConfigLogging describes logger configuration.
type ConfigLogging struct {
	LogLevel string `yaml:"loglevel"`
}

// ConfigPastes describes pastes subsystem configuration.
type ConfigPastes struct {
	Pagination int `yaml:"pagination"`
}
