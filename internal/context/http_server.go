package context

import (
	// local
	"gitlab.com/pztrn/fastpastebin/assets/static"

	// other
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (c *Context) initializeHTTPServer() {
	c.Echo = echo.New()
	c.Echo.Use(c.echoReqLogger())
	c.Echo.Use(middleware.Recover())
	c.Echo.DisableHTTP2 = true
	c.Echo.HideBanner = true
	c.Echo.HidePort = true

	// Static files.
	c.Echo.GET("/static/*", echo.WrapHandler(static.Handler))

	listenAddress := c.Config.HTTP.Address + ":" + c.Config.HTTP.Port

	go func() {
		c.Echo.Logger.Fatal(c.Echo.Start(listenAddress))
	}()
	c.Logger.Info().Msgf("Started HTTP server at %s", listenAddress)
}

// Wrapper around previous function.
func (c *Context) echoReqLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			c.Logger.Info().
				Str("IP", ec.RealIP()).
				Str("Host", ec.Request().Host).
				Str("Method", ec.Request().Method).
				Str("Path", ec.Request().URL.Path).
				Str("UA", ec.Request().UserAgent()).
				Msg("HTTP request")

			next(ec)
			return nil
		}
	}
}
