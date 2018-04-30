package http

import (
	// local
	"github.com/pztrn/fastpastebin/api/http/static"
	"github.com/pztrn/fastpastebin/context"

	// other
	"github.com/labstack/echo"
)

var (
	c *context.Context
)

// New initializes basic HTTP API, which shows only index page and serves
// static files.
func New(cc *context.Context) {
	c = cc
	c.Logger.Info().Msg("Initializing HTTP API...")

	// Static files.
	c.Echo.GET("/static/*", echo.WrapHandler(static.Handler))

	// Index.
	c.Echo.GET("/", indexGet)
}
