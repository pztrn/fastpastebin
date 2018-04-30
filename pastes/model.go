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
	// stdlib
	"time"

	// other
	"github.com/jmoiron/sqlx"
)

const (
	PASTE_KEEP_FOR_MINUTES = 1
	PASTE_KEEP_FOR_HOURS   = 2
	PASTE_KEEP_FOR_DAYS    = 3
	PASTE_KEEP_FOR_MONTHS  = 4
)

var (
	PASTE_KEEPS_CORELLATION = map[string]int{
		"M": PASTE_KEEP_FOR_MINUTES,
		"h": PASTE_KEEP_FOR_HOURS,
		"d": PASTE_KEEP_FOR_DAYS,
		"m": PASTE_KEEP_FOR_MONTHS,
	}
)

// Paste represents paste itself.
type Paste struct {
	ID              int        `db:"id"`
	Title           string     `db:"title"`
	Data            string     `db:"data"`
	CreatedAt       *time.Time `db:"created_at"`
	KeepFor         int        `db:"keep_for"`
	KeepForUnitType int        `db:"keep_for_unit_type"`
}

func (p *Paste) GetByID(db *sqlx.DB) error {
	err := db.Get(p, db.Rebind("SELECT * FROM `pastes` WHERE id=?"), p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Paste) GetPagedPastes(db *sqlx.DB, page int) ([]Paste, error) {
	var pastes []Paste

	// Pagination - 30 pastes on page.
	var startPagination = 0
	if page > 1 {
		startPagination = (page - 1) * 30
	}

	err := db.Select(&pastes, db.Rebind("SELECT * FROM `pastes` ORDER BY id DESC LIMIT 30 OFFSET ?"), startPagination)
	if err != nil {
		return nil, err
	}

	return pastes, nil

}

func (p *Paste) Save(db *sqlx.DB) (int64, error) {
	result, err := db.NamedExec("INSERT INTO `pastes` (title, data, created_at, keep_for, keep_for_unit_type) VALUES (:title, :data, :created_at, :keep_for, :keep_for_unit_type)", p)
	if err != nil {
		return 0, err
	}

	ID, err1 := result.LastInsertId()
	if err1 != nil {
		return 0, err
	}

	return ID, nil
}
