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

func New(cc *context.Context) {
	c = cc
	c.Logger.Info().Msg("Initializing HTTP API...")

	// Static files.
	c.Echo.GET("/static/*", echo.WrapHandler(static.Handler))

	// Index.
	c.Echo.GET("/", indexGet)

	// New paste.
	c.Echo.POST("/paste/", pastePOST)

	// Show paste.
	c.Echo.GET("/paste/:id", pasteGET)

	// Pastes list.
	c.Echo.GET("/pastes/", pastesGET)
	c.Echo.GET("/pastes/:page", pastesGET)
}
