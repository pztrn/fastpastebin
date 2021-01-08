// Fast Paste Bin - uberfast and easy-to-use pastebin.
//
// Copyright (c) 2018, Stanislav N. aka pztrn and Fast Paste Bin
// developers.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject
// to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package context

import (
	// stdlib
	"io/ioutil"
	"os"
	"path/filepath"

	// local
	"go.dev.pztrn.name/fastpastebin/internal/config"
	databaseinterface "go.dev.pztrn.name/fastpastebin/internal/database/interface"

	// other
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"go.dev.pztrn.name/flagger"
	"gopkg.in/yaml.v2"
)

// Context is a some sort of singleton. Basically it's a structure that
// initialized once and then passed to all parts of application. It
// contains everything every part of application need, like configuration
// access, logger, etc.
type Context struct {
	Config   *config.Struct
	Database databaseinterface.Interface
	Echo     *echo.Echo
	Flagger  *flagger.Flagger
	Logger   zerolog.Logger
}

// Initialize initializes context.
func (c *Context) Initialize() {
	c.initializeLogger()

	c.Flagger = flagger.New("fastpastebin", nil)
	c.Flagger.Initialize()

	_ = c.Flagger.AddFlag(&flagger.Flag{
		Name:         "config",
		Description:  "Configuration file path. Can be overridded with FASTPASTEBIN_CONFIG environment variable (this is what used in tests).",
		Type:         "string",
		DefaultValue: "NO_CONFIG",
	})
}

// InitializePost initializes everything that needs a configuration.
func (c *Context) InitializePost() {
	c.initializeLoggerPost()
	c.initializeHTTPServer()
}

// LoadConfiguration loads configuration and executes right after Flagger
// have parsed CLI flags, because it depends on "-config" defined in
// Initialize().
func (c *Context) LoadConfiguration() {
	c.Logger.Info().Msg("Loading configuration...")

	configPath := ""

	// We're accepting configuration path from "-config" CLI parameter
	// and FASTPASTEBIN_CONFIG environment variable. Later have higher
	// weight and can override "-config" value.
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
		c.Logger.Fatal().Err(err1).Msg("Failed to normalize path to configuration file")
	}

	c.Logger.Debug().Str("path", configPath).Msg("Configuration file path")

	c.Config = &config.Struct{}

	// Read configuration file.
	fileData, err2 := ioutil.ReadFile(normalizedConfigPath)
	if err2 != nil {
		c.Logger.Panic().Err(err2).Msg("Failed to read configuration file")
	}

	// Parse it into structure.
	err3 := yaml.Unmarshal(fileData, c.Config)
	if err3 != nil {
		c.Logger.Panic().Err(err3).Msg("Failed to parse configuration file")
	}

	// Yay! See what it gets!
	c.Logger.Debug().Msgf("Parsed configuration: %+v", c.Config)
}

// RegisterDatabaseInterface registers database interface for later use.
func (c *Context) RegisterDatabaseInterface(di databaseinterface.Interface) {
	c.Database = di
}

// RegisterEcho registers Echo instance for later usage.
func (c *Context) RegisterEcho(e *echo.Echo) {
	c.Echo = e
}

// Shutdown shutdowns entire application.
func (c *Context) Shutdown() {
	c.Logger.Info().Msg("Shutting down Fast Pastebin...")

	c.Database.Shutdown()
}
