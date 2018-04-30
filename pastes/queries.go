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

const (
	// Pagination. Hardcoded for 30 for now.
	PAGINATION = 10
)

// GetByID returns a single paste by ID.
func GetByID(id int) (*Paste, error) {
	p := &Paste{}
	dbConn := c.Database.GetDatabaseConnection()
	err := dbConn.Get(p, dbConn.Rebind("SELECT * FROM `pastes` WHERE id=?"), id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetPagedPastes returns a paged slice of pastes.
func GetPagedPastes(page int) ([]Paste, error) {
	var pastes []Paste
	dbConn := c.Database.GetDatabaseConnection()

	// Pagination - 30 pastes on page.
	var startPagination = 0
	if page > 1 {
		startPagination = (page - 1) * PAGINATION
	}

	err := dbConn.Select(&pastes, dbConn.Rebind("SELECT * FROM `pastes` ORDER BY id DESC LIMIT ? OFFSET ?"), PAGINATION, startPagination)
	if err != nil {
		return nil, err
	}

	return pastes, nil

}

// GetPastesPages returns an integer that represents quantity of pages
// that can be requested (or drawn in paginator).
func GetPastesPages() int {
	var pastesCount int
	dbConn := c.Database.GetDatabaseConnection()
	err := dbConn.Get(&pastesCount, "SELECT COUNT(id) FROM `pastes`")
	if err != nil {
		return 1
	}

	// Calculate pages.
	pages := pastesCount / PAGINATION
	// Check if we have any remainder. Add 1 to pages count if so.
	if pastesCount%PAGINATION != 0 {
		pages += 1
	}

	return pages
}

// Save saves paste to database and returns it's ID.
func Save(p *Paste) (int64, error) {
	dbConn := c.Database.GetDatabaseConnection()
	result, err := dbConn.NamedExec("INSERT INTO `pastes` (title, data, created_at, keep_for, keep_for_unit_type, language) VALUES (:title, :data, :created_at, :keep_for, :keep_for_unit_type, :language)", p)
	if err != nil {
		return 0, err
	}

	ID, err1 := result.LastInsertId()
	if err1 != nil {
		return 0, err
	}

	return ID, nil
}
