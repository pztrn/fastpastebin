package config

type ConfigStruct struct {
	Database ConfigDatabase `yaml:"database"`
	Logging  ConfigLogging  `yaml:"logging"`
	HTTP     ConfigHTTP     `yaml:"http"`
}
