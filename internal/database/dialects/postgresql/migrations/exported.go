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
	// local
	"go.dev.pztrn.name/fastpastebin/internal/context"

	// other
	//"gitlab.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

var (
	c *context.Context
)

// New initializes migrations.
func New(cc *context.Context) {
	c = cc
}

// Migrate launching migrations.
func Migrate() {
	c.Logger.Info().Msg("Migrating database...")

	goose.SetDialect("postgres")
	goose.AddNamedMigration("1_initial.go", InitialUp, nil)
	goose.AddNamedMigration("2_paste_lang.go", PasteLangUp, PasteLangDown)
	goose.AddNamedMigration("3_private_pastes.go", PrivatePastesUp, PrivatePastesDown)
	goose.AddNamedMigration("4_passworded_pastes.go", PasswordedPastesUp, PasswordedPastesDown)
	// Add new migrations BEFORE this message.

	dbConn := c.Database.GetDatabaseConnection()
	if dbConn != nil {
		err := goose.Up(dbConn, ".")
		if err != nil {
			c.Logger.Info().Msgf("%+v", err)
			c.Logger.Panic().Msgf("Failed to migrate database to latest version: %s", err.Error())
		}
	} else {
		c.Logger.Warn().Msg("Current database dialect isn't supporting migrations, skipping")
	}
}
