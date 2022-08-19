package application

import "errors"

var (
	// ErrConfigurationError indicates that this error is related to configuration.
	ErrConfigurationError = errors.New("configuration")

	// ErrConfigurationLoad indicates that error appears when trying to load
	// configuration data from file.
	ErrConfigurationLoad = errors.New("loading configuration")

	// ErrConfigurationPathNotDefined indicates that CONFIG_PATH environment variable is empty or not defined.
	ErrConfigurationPathNotDefined = errors.New("configuration path (CONFIG_PATH) is empty or not defined")
)
