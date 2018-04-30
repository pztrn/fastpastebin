package api

import (
	// local
	"github.com/pztrn/fastpastebin/api/http"
	"github.com/pztrn/fastpastebin/api/json"
	"github.com/pztrn/fastpastebin/context"

	// other
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	c *context.Context
	e *echo.Echo
)

func New(cc *context.Context) {
	c = cc
}

func InitializeAPI() {
	c.Logger.Info().Msg("Initializing HTTP server...")

	e = echo.New()
	e.Use(echoReqLogger())
	e.Use(middleware.Recover())
	e.DisableHTTP2 = true
	e.HideBanner = true
	e.HidePort = true

	c.RegisterEcho(e)

	http.New(c)
	json.New(c)

	listenAddress := c.Config.HTTP.Address + ":" + c.Config.HTTP.Port

	go func() {
		e.Logger.Fatal(e.Start(listenAddress))
	}()
	c.Logger.Info().Msgf("Started HTTP server at %s", listenAddress)
}
