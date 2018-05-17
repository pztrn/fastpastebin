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
	//"github.com/alecthomas/chroma"
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
	Language        string     `db:"language"`
	Private         bool       `db:"private"`
	Password        string     `db:"password"`
	PasswordSalt    string     `db:"password_salt"`
}

func (p *Paste) GetExpirationTime() time.Time {
	var expirationTime time.Time
	switch p.KeepForUnitType {
	case PASTE_KEEP_FOR_MINUTES:
		expirationTime = p.CreatedAt.Add(time.Minute * time.Duration(p.KeepFor))
	case PASTE_KEEP_FOR_HOURS:
		expirationTime = p.CreatedAt.Add(time.Hour * time.Duration(p.KeepFor))
	case PASTE_KEEP_FOR_DAYS:
		expirationTime = p.CreatedAt.Add(time.Hour * 24 * time.Duration(p.KeepFor))
	case PASTE_KEEP_FOR_MONTHS:
		expirationTime = p.CreatedAt.Add(time.Hour * 24 * 30 * time.Duration(p.KeepFor))
	}

	return expirationTime
}

// IsExpired checks if paste is already expired (or not).
func (p *Paste) IsExpired() bool {
	curTime := time.Now().UTC()
	expirationTime := p.GetExpirationTime()

	if curTime.Sub(expirationTime).Seconds() > 0 {
		return true
	}

	return false
}
