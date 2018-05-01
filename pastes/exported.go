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

package pastes

import (
	// local
	"github.com/pztrn/fastpastebin/context"
)

var (
	c *context.Context
)

// New initializes pastes package and adds neccessary HTTP and API
// endpoints.
func New(cc *context.Context) {
	c = cc

	// New paste.
	c.Echo.POST("/paste/", pastePOST)

	// Show public paste.
	c.Echo.GET("/paste/:id", pasteGET)
	// Show RAW representation of public paste.
	c.Echo.GET("/paste/:id/raw", pasteRawGET)

	// Show private paste.
	c.Echo.GET("/paste/:id/:timestamp", pasteGET)
	// Show RAW representation of private paste.
	c.Echo.GET("/paste/:id/:timestamp/raw", pasteRawGET)

	// Pastes list.
	c.Echo.GET("/pastes/", pastesGET)
	c.Echo.GET("/pastes/:page", pastesGET)
}
