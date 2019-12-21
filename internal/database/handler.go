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

package database

import (
	// stdlib
	"database/sql"

	// local

	dialectinterface "go.dev.pztrn.name/fastpastebin/internal/database/dialects/interface"
	"go.dev.pztrn.name/fastpastebin/internal/structs"
)

// Handler is an interfaceable structure that proxifies calls from anyone
// to Database structure.
type Handler struct{}

func (dbh Handler) GetDatabaseConnection() *sql.DB {
	return d.GetDatabaseConnection()
}

func (dbh Handler) GetPaste(pasteID int) (*structs.Paste, error) {
	return d.GetPaste(pasteID)
}

func (dbh Handler) GetPagedPastes(page int) ([]structs.Paste, error) {
	return d.GetPagedPastes(page)
}

func (dbh Handler) GetPastesPages() int {
	return d.GetPastesPages()
}

// Initialize initializes connection to database.
func (dbh Handler) Initialize() {
	d.Initialize()
}

func (dbh Handler) RegisterDialect(di dialectinterface.Interface) {
	d.RegisterDialect(di)
}

func (dbh Handler) SavePaste(p *structs.Paste) (int64, error) {
	return d.SavePaste(p)
}

func (dbh Handler) Shutdown() {
	d.Shutdown()
}
