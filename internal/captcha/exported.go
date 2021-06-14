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

package captcha

import (
	"github.com/dchest/captcha"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"go.dev.pztrn.name/fastpastebin/internal/context"
)

var (
	c   *context.Context
	log zerolog.Logger
)

// New initializes captcha package and adds necessary HTTP and API
// endpoints.
func New(cc *context.Context) {
	c = cc
	log = c.Logger.With().Str("type", "internal").Str("package", "captcha").Logger()

	// New paste.
	c.Echo.GET("/captcha/:id.png", echo.WrapHandler(captcha.Server(captcha.StdWidth, captcha.StdHeight)))
}

// NewCaptcha creates new captcha string.
func NewCaptcha() string {
	s := captcha.New()
	log.Debug().Str("captcha string", s).Msg("Created new captcha string")

	return s
}

// Verify verifies captchas.
func Verify(captchaString string, captchaSolution string) bool {
	return captcha.VerifyString(captchaString, captchaSolution)
}
