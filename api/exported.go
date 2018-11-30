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
	// local
	"gitlab.com/pztrn/fastpastebin/api/http"
	"gitlab.com/pztrn/fastpastebin/api/json"
	"gitlab.com/pztrn/fastpastebin/context"

	// other
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	c *context.Context
	e *echo.Echo
)

// New initializes variables for api package.
func New(cc *context.Context) {
	c = cc
}

// InitializeAPI initializes HTTP API and starts web server.
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
