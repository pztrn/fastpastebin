package config

// ConfigLogging describes logger configuration.
type ConfigLogging struct {
	LogToFile bool   `yaml:"log_to_file"`
	FileName  string `yaml:"filename"`
	LogLevel  string `yaml:"loglevel"`
}
