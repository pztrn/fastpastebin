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
	// stdlib
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	// other
	"golang.org/x/crypto/scrypt"
)

const (
	PasteKeepForever    = 0
	PasteKeepForMinutes = 1
	PasteKeepForHours   = 2
	PasteKeepForDays    = 3
	PasteKeepForMonths  = 4

	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	PasteKeepsCorellation = map[string]int{
		"M":       PasteKeepForMinutes,
		"h":       PasteKeepForHours,
		"d":       PasteKeepForDays,
		"m":       PasteKeepForMonths,
		"forever": PasteKeepForever,
	}
)

// Paste represents paste itself.
type Paste struct {
	ID              int        `db:"id" json:"id"`
	Title           string     `db:"title" json:"title"`
	Data            string     `db:"data" json:"data"`
	CreatedAt       *time.Time `db:"created_at" json:"created_at"`
	KeepFor         int        `db:"keep_for" json:"keep_for"`
	KeepForUnitType int        `db:"keep_for_unit_type" json:"keep_for_unit_type"`
	Language        string     `db:"language" json:"language"`
	Private         bool       `db:"private" json:"private"`
	Password        string     `db:"password" json:"password"`
	PasswordSalt    string     `db:"password_salt" json:"password_salt"`
}

// CreatePassword creates password for current paste.
func (p *Paste) CreatePassword(password string) error {
	// Create salt - random string.
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
