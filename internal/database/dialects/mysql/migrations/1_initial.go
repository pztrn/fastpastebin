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

package migrations

import (
	"database/sql"
)

func InitialUp(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE pastes (
		id int(11) NOT NULL AUTO_INCREMENT COMMENT 'Paste ID', 
		title text NOT NULL COMMENT 'Paste title', 
		data longtext NOT NULL COMMENT 'Paste data', 
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Paste creation timestamp', 
		keep_for int(4) NOT NULL DEFAULT 1 COMMENT 'Keep for integer. 0 - forever.', 
		keep_for_unit_type int(1) NOT NULL DEFAULT 1 COMMENT 'Keep for unit type. 1 - minutes, 2 - hours, 3 - days, 4 - months.',
		PRIMARY KEY (id), UNIQUE KEY id (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='Pastes';`)
	if err != nil {
		// nolint:wrapcheck
		return err
	}

	return nil
}
