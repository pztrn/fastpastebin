package context

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.dev.pztrn.name/fastpastebin/assets"
)

func (c *Context) initializeHTTPServer() {
	c.Echo = echo.New()
	c.Echo.Use(c.echoReqLogger())
	c.Echo.Use(middleware.Recover())
	c.Echo.Use(middleware.BodyLimit(c.Config.HTTP.MaxBodySizeMegabytes + "M"))
	c.Echo.DisableHTTP2 = true
	c.Echo.HideBanner = true
	c.Echo.HidePort = true

	// Static files.
	c.Echo.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.Data))))

	listenAddress := c.Config.HTTP.Address + ":" + c.Config.HTTP.Port

	go func() {
		c.Echo.Logger.Fatal(c.Echo.Start(listenAddress))
	}()
	c.Logger.Info().Str("address", listenAddress).Msg("Started HTTP server")
}

// Wrapper around previous function.
func (c *Context) echoReqLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			c.Logger.Info().
				Str("IP", ectx.RealIP()).
				Str("Host", ectx.Request().Host).
				Str("Method", ectx.Request().Method).
				Str("Path", ectx.Request().URL.Path).
				Str("UA", ectx.Request().UserAgent()).
				Msg("HTTP request")

			return next(ectx)
		}
	}
}
