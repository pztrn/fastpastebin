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

package structs

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/scrypt"
)

const (
	// PasteKeepForever indicates that paste should be kept forever.
	PasteKeepForever = 0
	// PasteKeepForMinutes indicates that saved timeout is in minutes.
	PasteKeepForMinutes = 1
	// PasteKeepForHours indicates that saved timeout is in hours.
	PasteKeepForHours = 2
	// PasteKeepForDays indicates that saved timeout is in days.
	PasteKeepForDays = 3
	// PasteKeepForMonths indicates that saved timeout is in months.
	PasteKeepForMonths = 4

	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// PasteKeepsCorrelation is a correlation map between database representation
// and passed data representation.
var PasteKeepsCorrelation = map[string]int{
	"M":       PasteKeepForMinutes,
	"h":       PasteKeepForHours,
	"d":       PasteKeepForDays,
	"m":       PasteKeepForMonths,
	"forever": PasteKeepForever,
}

// Paste represents paste itself.
type Paste struct {
	CreatedAt       *time.Time `db:"created_at" json:"created_at"`
	Title           string     `db:"title" json:"title"`
	Data            string     `db:"data" json:"data"`
	Language        string     `db:"language" json:"language"`
	Password        string     `db:"password" json:"password"`
	PasswordSalt    string     `db:"password_salt" json:"password_salt"`
	ID              int        `db:"id" json:"id"`
	KeepFor         int        `db:"keep_for" json:"keep_for"`
	KeepForUnitType int        `db:"keep_for_unit_type" json:"keep_for_unit_type"`
	Private         bool       `db:"private" json:"private"`
}

// CreatePassword creates password for current paste.
func (p *Paste) CreatePassword(password string) error {
	// Create salt - random string.
	// Yes, it is insecure. Should be refactored!
	// nolint:gosec
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	saltBytes := make([]byte, 64)

	for i := range saltBytes {
		saltBytes[i] = charset[seededRand.Intn(len(charset))]
	}

	saltHashBytes := sha256.Sum256(saltBytes)
	p.PasswordSalt = fmt.Sprintf("%x", saltHashBytes)

	// Create crypted password and hash it.
	passwordCrypted, err := scrypt.Key([]byte(password), []byte(p.PasswordSalt), 131072, 8, 1, 64)
	if err != nil {
		// nolint:wrapcheck
		return err
	}

	passwordHashBytes := sha256.Sum256(passwordCrypted)
	p.Password = fmt.Sprintf("%x", passwordHashBytes)

	return nil
}

// GenerateCryptedCookieValue generates crypted cookie value for paste.
func (p *Paste) GenerateCryptedCookieValue() string {
	cookieValueCrypted, _ := scrypt.Key([]byte(p.Password), []byte(p.PasswordSalt), 131072, 8, 1, 64)

	return fmt.Sprintf("%x", sha256.Sum256(cookieValueCrypted))
}

func (p *Paste) GetExpirationTime() time.Time {
	var expirationTime time.Time

	switch p.KeepForUnitType {
	case PasteKeepForever:
		expirationTime = time.Now().UTC().Add(time.Hour * 1)
	case PasteKeepForMinutes:
		expirationTime = p.CreatedAt.Add(time.Minute * time.Duration(p.KeepFor))
	case PasteKeepForHours:
		expirationTime = p.CreatedAt.Add(time.Hour * time.Duration(p.KeepFor))
	case PasteKeepForDays:
		expirationTime = p.CreatedAt.Add(time.Hour * 24 * time.Duration(p.KeepFor))
	case PasteKeepForMonths:
		expirationTime = p.CreatedAt.Add(time.Hour * 24 * 30 * time.Duration(p.KeepFor))
	}

	return expirationTime
}

// IsExpired checks if paste is already expired (or not).
func (p *Paste) IsExpired() bool {
	curTime := time.Now().UTC()
	expirationTime := p.GetExpirationTime()

	return curTime.Sub(expirationTime).Seconds() > 0
}

// VerifyPassword verifies that provided password is valid.
func (p *Paste) VerifyPassword(password string) bool {
	// Create crypted password and hash it.
	passwordCrypted, err := scrypt.Key([]byte(password), []byte(p.PasswordSalt), 131072, 8, 1, 64)
	if err != nil {
		return false
	}

	passwordHashBytes := sha256.Sum256(passwordCrypted)
	providedPassword := fmt.Sprintf("%x", passwordHashBytes)

	return providedPassword == p.Password
}
