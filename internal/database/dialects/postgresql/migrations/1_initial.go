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
	_, err := tx.Exec(`
					CREATE TABLE pastes 
					(
						id 					SERIAL PRIMARY KEY,
						title 				TEXT NOT NULL,
						data 				TEXT NOT NULL,
						created_at 			TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
						keep_for 			INTEGER NOT NULL DEFAULT 1,
						keep_for_unit_type 	SMALLINT NOT NULL DEFAULT 1
					);

					COMMENT ON COLUMN pastes.id IS 'Paste ID';
					COMMENT ON COLUMN pastes.title IS 'Paste title';
					COMMENT ON COLUMN pastes.data IS 'Paste data';
					COMMENT ON COLUMN pastes.created_at IS 'Paste creation timestamp';
					COMMENT ON COLUMN pastes.keep_for IS 'Keep for integer. 0 - forever.';
					COMMENT ON COLUMN pastes.keep_for_unit_type IS 'Keep for unit type. 0 - forever, 1 - minutes, 2 - hours, 3 - days, 4 - months.';
	`)
	if err != nil {
		// nolint:wrapcheck
		return err
	}

	return nil
}
