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
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	event.Str("memalloc", fmt.Sprintf("%dMB", m.Alloc/1024/1024))
	event.Str("memsys", fmt.Sprintf("%dMB", m.Sys/1024/1024))
	event.Str("numgc", fmt.Sprintf("%d", m.NumGC))
}

// Initializes logger.
func (c *Context) initializeLogger() {
	// Устанавливаем форматирование логгера.
	// nolint:exhaustivestruct
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.RFC3339}
	output.FormatLevel = func(lvlRaw interface{}) string {
		var v string

		if lvl, ok := lvlRaw.(string); ok {
			lvl = strings.ToUpper(lvl)
			switch lvl {
			case "DEBUG":
				v = fmt.Sprintf("\x1b[30m%-5s\x1b[0m", lvl)
			case "ERROR":
				v = fmt.Sprintf("\x1b[31m%-5s\x1b[0m", lvl)
			case "FATAL":
				v = fmt.Sprintf("\x1b[35m%-5s\x1b[0m", lvl)
			case "INFO":
				v = fmt.Sprintf("\x1b[32m%-5s\x1b[0m", lvl)
			case "PANIC":
				v = fmt.Sprintf("\x1b[36m%-5s\x1b[0m", lvl)
			case "WARN":
				v = fmt.Sprintf("\x1b[33m%-5s\x1b[0m", lvl)
			default:
				v = lvl
			}
		}

		return fmt.Sprintf("| %s |", v)
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
