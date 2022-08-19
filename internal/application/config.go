package application

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"go.dev.pztrn.name/fastpastebin/internal/helpers"
	"gopkg.in/yaml.v2"
)

// Config represents configuration structure.
type Config struct {
	app *Application
	log zerolog.Logger

	Database ConfigDatabase `yaml:"database"`
	Logging  ConfigLogging  `yaml:"logging"`
	HTTP     ConfigHTTP     `yaml:"http"`
	Pastes   ConfigPastes   `yaml:"pastes"`
}

func newConfig(app *Application) (*Config, error) {
	//nolint:exhaustruct
	cfg := &Config{
		app: app,
		log: app.Log.With().Str("type", "core").Str("name", "configuration").Logger(),
	}

	if err := cfg.initialize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrConfigurationError, err)
	}

	return cfg, nil
}

func (c *Config) initialize() error {
	c.log.Info().Msg("Initializing configuration...")

	configPathRaw, found := os.LookupEnv("FASTPASTEBIN_CONFIG")
	if !found {
		return fmt.Errorf("%s: %w", ErrConfigurationLoad, ErrConfigurationPathNotDefined)
	}

	configPath, err := helpers.NormalizePath(configPathRaw)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrConfigurationLoad, err)
	}

	c.log.Info().Str("config path", configPath).Msg("Reading configuration file...")

	fileData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrConfigurationLoad, err)
	}

	if err := yaml.Unmarshal(fileData, c); err != nil {
		return fmt.Errorf("%s: %w", ErrConfigurationLoad, err)
	}

	c.log.Debug().Msgf("Configuration loaded: %+v", c)

	return nil
}
