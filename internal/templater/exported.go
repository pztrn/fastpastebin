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

package templater

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"go.dev.pztrn.name/fastpastebin/assets"
	"go.dev.pztrn.name/fastpastebin/internal/context"
)

var (
	c   *context.Context
	log zerolog.Logger
)

// GetErrorTemplate returns formatted error template.
// If error.html wasn't found - it will return "error.html not found"
// message as simple string.
func GetErrorTemplate(ec echo.Context, errorText string) string {
	// Getting main error template.
	mainhtml := GetTemplate(ec, "error.html", map[string]string{"error": errorText})

	return mainhtml
}

// GetRawTemplate returns only raw template data.
func GetRawTemplate(ec echo.Context, templateName string, data map[string]string) string {
	// Getting main template.
	tplRaw, err := assets.Data.ReadFile(templateName)
	if err != nil {
		_ = ec.String(http.StatusBadRequest, templateName+" not found.")

		return ""
	}

	tpl := string(tplRaw)
	// Replace placeholders with data from data map.
	for placeholder, value := range data {
		tpl = strings.Replace(tpl, "{"+placeholder+"}", value, -1)
	}

	return tpl
}

// GetTemplate returns formatted template that can be outputted to client.
func GetTemplate(ec echo.Context, name string, data map[string]string) string {
	log.Debug().Str("name", name).Msg("Requested template")

	// Getting main template.
	mainhtml, err := assets.Data.ReadFile("main.html")
	if err != nil {
		_ = ec.String(http.StatusBadRequest, "main.html not found.")

		return ""
	}

	// Getting navigation.
	navhtml, err1 := assets.Data.ReadFile("navigation.html")
	if err1 != nil {
		_ = ec.String(http.StatusBadRequest, "navigation.html not found.")

		return ""
	}

	// Getting footer.
	footerhtml, err2 := assets.Data.ReadFile("footer.html")
	if err2 != nil {
		_ = ec.String(http.StatusBadRequest, "footer.html not found.")

		return ""
	}

	// Format main template.
	tpl := strings.Replace(string(mainhtml), "{navigation}", string(navhtml), 1)
	tpl = strings.Replace(tpl, "{footer}", string(footerhtml), 1)
	// Version.
	tpl = strings.Replace(tpl, "{version}", context.Version, 1)

	// Get requested template.
	reqhtml, err3 := assets.Data.ReadFile(name)
	if err3 != nil {
		_ = ec.String(http.StatusBadRequest, name+" not found.")

		return ""
	}

	// Replace documentBody.
	tpl = strings.Replace(tpl, "{documentBody}", string(reqhtml), 1)

	// Replace placeholders with data from data map.
	for placeholder, value := range data {
		tpl = strings.Replace(tpl, "{"+placeholder+"}", value, -1)
	}

	return tpl
}

// Initialize initializes package.
func Initialize(cc *context.Context) {
	c = cc
	log = c.Logger.With().Str("type", "internal").Str("package", "templater").Logger()
}
