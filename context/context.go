package context

import (
	// stdlib
	"io/ioutil"
	"os"
	"path/filepath"

	// local
	"github.com/pztrn/fastpastebin/config"
	"github.com/pztrn/fastpastebin/database/interface"

	// other
	"github.com/labstack/echo"
	"github.com/pztrn/flagger"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type Context struct {
	Config   *config.ConfigStruct
	Database databaseinterface.Interface
	Echo     *echo.Echo
	Flagger  *flagger.Flagger
	Logger   zerolog.Logger
}

func (c *Context) Initialize() {
	c.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Caller().Logger()

	c.Flagger = flagger.New(nil)
	c.Flagger.Initialize()

	c.Flagger.AddFlag(&flagger.Flag{
		Name:         "config",
		Description:  "Configuration file path. Can be overridded with FASTPASTEBIN_CONFIG environment variable (this is what used in tests).",
		Type:         "string",
		DefaultValue: "NO_CONFIG",
	})
}

func (c *Context) LoadConfiguration() {
	c.Logger.Info().Msg("Loading configuration...")

	var configPath = ""

	configPathFromCLI, err := c.Flagger.GetStringValue("config")
	configPathFromEnv, configPathFromEnvFound := os.LookupEnv("FASTPASTEBIN_CONFIG")

	if err != nil && configPathFromEnvFound || err == nil && configPathFromEnvFound {
		configPath = configPathFromEnv
	} else if err != nil && !configPathFromEnvFound || err == nil && configPathFromCLI == "NO_CONFIG" {
		c.Logger.Panic().Msg("Configuration file path wasn't passed via '-config' or 'FASTPASTEBIN_CONFIG' environment variable. Cannot continue.")
	} else if err == nil && !configPathFromEnvFound {
		configPath = configPathFromCLI
	}

	// Normalize file path.
	normalizedConfigPath, err1 := filepath.Abs(configPath)
	if err1 != nil {
		c.Logger.Fatal().Msgf("Failed to normalize path to configuration file: %s", err1.Error())
	}

	c.Logger.Debug().Msgf("Configuration file path: %s", configPath)

	c.Config = &config.ConfigStruct{}

	fileData, err2 := ioutil.ReadFile(normalizedConfigPath)
	if err2 != nil {
		c.Logger.Panic().Msgf("Failed to read configuration file: %s", err2.Error())
	}

	err3 := yaml.Unmarshal(fileData, c.Config)
	if err3 != nil {
		c.Logger.Panic().Msgf("Failed to parse configuration file: %s", err3.Error())
	}

	c.Logger.Debug().Msgf("Parsed configuration: %+v", c.Config)
}

func (c *Context) RegisterDatabaseInterface(di databaseinterface.Interface) {
	c.Database = di
}

func (c *Context) RegisterEcho(e *echo.Echo) {
	c.Echo = e
}

func (c *Context) Shutdown() {
	c.Logger.Info().Msg("Shutting down Fast Pastebin...")
}
