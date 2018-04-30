// Fast Paste Bin - uberfast and easy-to-use pastebin.
//
// Copyright (c) 2018, Stanislav N. aka pztrn and Fast Paste Bin
// developers.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject
// to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
