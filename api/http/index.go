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

package http

import (
	// stdlib
	"net/http"
	"strings"

	// local
	"github.com/pztrn/fastpastebin/api/http/static"

	// other
	"github.com/alecthomas/chroma/lexers"
	"github.com/labstack/echo"
)

// Index of this site.
func indexGet(ec echo.Context) error {
	htmlRaw, err := static.ReadFile("index.html")
	if err != nil {
		return ec.String(http.StatusNotFound, "index.html wasn't found!")
	}

	// Generate list of available languages to highlight.
	availableLexers := lexers.Names(false)

	var availableLexersSelectOpts = "<option value='text'>Text</option><option value='autodetect'>Auto detect</option><option disabled>-----</option>"
	for i := range availableLexers {
		availableLexersSelectOpts += "<option value='" + availableLexers[i] + "'>" + availableLexers[i] + "</option>"
	}

	html := strings.Replace(string(htmlRaw), "{lexers}", availableLexersSelectOpts, 1)

	return ec.HTML(http.StatusOK, html)
}
