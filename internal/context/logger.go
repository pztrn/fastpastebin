package context

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Puts memory usage into log lines.
func (c *Context) getMemoryUsage(event *zerolog.Event, level zerolog.Level, message string) {
	var memstats runtime.MemStats

	runtime.ReadMemStats(&memstats)

	event.Str("memalloc", fmt.Sprintf("%dMB", memstats.Alloc/1024/1024))
	event.Str("memsys", fmt.Sprintf("%dMB", memstats.Sys/1024/1024))
	event.Str("numgc", fmt.Sprintf("%d", memstats.NumGC))
}

// Initializes logger.
func (c *Context) initializeLogger() {
	// Устанавливаем форматирование логгера.
	//nolint:exhaustruct
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.RFC3339}
	output.FormatLevel = func(lvlRaw interface{}) string {
		var lvl string

		if lvlAsString, ok := lvlRaw.(string); ok {
			lvlAsString = strings.ToUpper(lvlAsString)
			switch lvlAsString {
			case "DEBUG":
				lvl = fmt.Sprintf("\x1b[30m%-5s\x1b[0m", lvlAsString)
			case "ERROR":
				lvl = fmt.Sprintf("\x1b[31m%-5s\x1b[0m", lvlAsString)
			case "FATAL":
				lvl = fmt.Sprintf("\x1b[35m%-5s\x1b[0m", lvlAsString)
			case "INFO":
				lvl = fmt.Sprintf("\x1b[32m%-5s\x1b[0m", lvlAsString)
			case "PANIC":
				lvl = fmt.Sprintf("\x1b[36m%-5s\x1b[0m", lvlAsString)
			case "WARN":
				lvl = fmt.Sprintf("\x1b[33m%-5s\x1b[0m", lvlAsString)
			default:
				lvl = lvlAsString
			}
		}

		return fmt.Sprintf("| %s |", lvl)
	}

	c.Logger = zerolog.New(output).With().Timestamp().Logger()

	c.Logger = c.Logger.Hook(zerolog.HookFunc(c.getMemoryUsage))
}

// Initialize logger after configuration parse.
func (c *Context) initializeLoggerPost() {
	// Set log level.
	c.Logger.Info().Msgf("Setting logger level: %s", c.Config.Logging.LogLevel)

	switch c.Config.Logging.LogLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}
}
