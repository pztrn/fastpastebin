package api

import (
	// other
	"github.com/labstack/echo"
)

// Logs Echo requests.
func echoReqLog(ec echo.Context, next echo.HandlerFunc) error {
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

// Wrapper around previous function.
func echoReqLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return echoReqLog(c, next)
		}
	}
}
