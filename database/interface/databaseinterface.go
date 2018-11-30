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

package databaseinterface

import (
	// stdlib
	"database/sql"

	// local
	"gitlab.com/pztrn/fastpastebin/database/dialects/interface"
	"gitlab.com/pztrn/fastpastebin/pastes/model"
)

// Interface represents database interface which is available to all
// parts of application and registers with context.Context.
type Interface interface {
	GetDatabaseConnection() *sql.DB
	GetPaste(pasteID int) (*pastesmodel.Paste, error)
	GetPagedPastes(page int) ([]pastesmodel.Paste, error)
	GetPastesPages() int
	Initialize()
	RegisterDialect(dialectinterface.Interface)
	SavePaste(p *pastesmodel.Paste) (int64, error)
	Shutdown()
}
