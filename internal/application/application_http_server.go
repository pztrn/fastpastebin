package application

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.dev.pztrn.name/fastpastebin/assets"
)

// Wrapper around previous function.
func (a *Application) echoReqLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			a.Log.Info().
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

func (a *Application) initializeHTTPServer() {
	a.Echo = echo.New()
	a.Echo.Use(a.echoReqLogger())
	a.Echo.Use(middleware.Recover())
	a.Echo.Use(middleware.BodyLimit(a.Config.HTTP.MaxBodySizeMegabytes + "M"))
	a.Echo.DisableHTTP2 = true
	a.Echo.HideBanner = true
	a.Echo.HidePort = true

	// Static files.
	a.Echo.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.Data))))
}

func (a *Application) startHTTPServer() {
	listenAddress := a.Config.HTTP.Address + ":" + a.Config.HTTP.Port

	go func() {
		a.Echo.Logger.Fatal(a.Echo.Start(listenAddress))
	}()
	a.Log.Info().Str("address", listenAddress).Msg("Started HTTP server")
}
