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

package indexpage

import (
	"net/http"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/labstack/echo"
	"go.dev.pztrn.name/fastpastebin/internal/captcha"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/flatfiles"
	"go.dev.pztrn.name/fastpastebin/internal/templater"
)

// Index of this site.
func indexGet(ectx echo.Context) error {
	// We should check if database connection available.
	dbConn := ctx.Database.GetDatabaseConnection()
	if ctx.Config.Database.Type != flatfiles.FlatFileDialect && dbConn == nil {
		//nolint:wrapcheck
		return ectx.Redirect(http.StatusFound, "/database_not_available")
	}

	// Generate list of available languages to highlight.
	availableLexers := lexers.Names(false)

	availableLexersSelectOpts := "<option value='text'>Text</option><option value='autodetect'>Auto detect</option><option disabled>-----</option>"
	for i := range availableLexers {
		availableLexersSelectOpts += "<option value='" + availableLexers[i] + "'>" + availableLexers[i] + "</option>"
	}

	// Captcha.
	captchaString := captcha.NewCaptcha()

	htmlData := templater.GetTemplate(ectx, "index.html", map[string]string{"lexers": availableLexersSelectOpts, "captchaString": captchaString})

	//nolint:wrapcheck
	return ectx.HTML(http.StatusOK, htmlData)
}
