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

func PrivatePastesUp(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE pastes ADD COLUMN private BOOLEAN NOT NULL DEFAULT false; COMMENT ON COLUMN pastes.private IS 'Private paste? If true - additional URL parameter (UNIX TIMESTAMP) of paste will be required to access.';")
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	return nil
}

func PrivatePastesDown(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE pastes DROP COLUMN private")
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	return nil
}
