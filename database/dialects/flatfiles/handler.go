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

package flatfiles

import (
	// stdlib
	"database/sql"

	// local
	"lab.pztrn.name/pztrn/fastpastebin/pastes/model"
)

type Handler struct{}

func (dbh Handler) GetDatabaseConnection() *sql.DB {
	return f.GetDatabaseConnection()
}

func (dbh Handler) GetPaste(pasteID int) (*pastesmodel.Paste, error) {
	return f.GetPaste(pasteID)
}

func (dbh Handler) GetPagedPastes(page int) ([]pastesmodel.Paste, error) {
	return f.GetPagedPastes(page)
}

func (dbh Handler) GetPastesPages() int {
	return f.GetPastesPages()
}

func (dbh Handler) Initialize() {
	f.Initialize()
}

func (dbh Handler) SavePaste(p *pastesmodel.Paste) (int64, error) {
	return f.SavePaste(p)
}

func (dbh Handler) Shutdown() {
	f.Shutdown()
}
