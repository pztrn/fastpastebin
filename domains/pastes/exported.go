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
	"regexp"

	"go.dev.pztrn.name/fastpastebin/internal/context"
)

var regexInts = regexp.MustCompile("[0-9]+")

var ctx *context.Context

// New initializes pastes package and adds necessary HTTP and API
// endpoints.
func New(cc *context.Context) {
	ctx = cc

	////////////////////////////////////////////////////////////
	// HTTP endpoints.
	////////////////////////////////////////////////////////////
	// New paste.
	ctx.Echo.POST("/paste/", pastePOSTWebInterface)
	// Show public paste.
	ctx.Echo.GET("/paste/:id", pasteGETWebInterface)
	// Show RAW representation of public paste.
	ctx.Echo.GET("/paste/:id/raw", pasteRawGETWebInterface)
	// Show private paste.
	ctx.Echo.GET("/paste/:id/:timestamp", pasteGETWebInterface)
	// Show RAW representation of private paste.
	ctx.Echo.GET("/paste/:id/:timestamp/raw", pasteRawGETWebInterface)
	// Verify access to passworded paste.
	ctx.Echo.GET("/paste/:id/:timestamp/verify", pastePasswordedVerifyGet)
	ctx.Echo.POST("/paste/:id/:timestamp/verify", pastePasswordedVerifyPost)
	// Pastes list.
	ctx.Echo.GET("/pastes/", pastesGET)
	ctx.Echo.GET("/pastes/:page", pastesGET)
}
