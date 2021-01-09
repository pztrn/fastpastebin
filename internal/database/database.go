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

package database

import (
	// stdlib
	"database/sql"
	"time"

	// local
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/flatfiles"
	dialectinterface "go.dev.pztrn.name/fastpastebin/internal/database/dialects/interface"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/mysql"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/postgresql"
	"go.dev.pztrn.name/fastpastebin/internal/structs"

	// other
	_ "github.com/go-sql-driver/mysql"
)

// Database represents control structure for database connection.
type Database struct {
	db dialectinterface.Interface
}

// Database cleanup function. Executes once per hour, hardcoded for now and is
// a subject of change in future.
func (db *Database) cleanup() {
	for {
		c.Logger.Info().Msg("Starting pastes cleanup procedure...")

		pages := db.db.GetPastesPages()

		var pasteIDsToRemove []int

		for i := 0; i < pages; i++ {
			pastes, err := db.db.GetPagedPastes(i)
			if err != nil {
				c.Logger.Error().Err(err).Int("page", i).Msg("Failed to perform database cleanup")
			}

			for _, paste := range pastes {
				if paste.IsExpired() {
					pasteIDsToRemove = append(pasteIDsToRemove, paste.ID)
				}
			}
		}

		for _, pasteID := range pasteIDsToRemove {
			err := db.DeletePaste(pasteID)
			if err != nil {
				c.Logger.Error().Err(err).Int("paste", pasteID).Msg("Failed to delete paste!")
			}
		}

		c.Logger.Info().Msg("Pastes cleanup done.")

		time.Sleep(time.Hour)
	}
}

func (db *Database) DeletePaste(pasteID int) error {
	return db.db.DeletePaste(pasteID)
}

func (db *Database) GetDatabaseConnection() *sql.DB {
	if db.db != nil {
		return db.db.GetDatabaseConnection()
	}

	return nil
}

func (db *Database) GetPaste(pasteID int) (*structs.Paste, error) {
	return db.db.GetPaste(pasteID)
}

func (db *Database) GetPagedPastes(page int) ([]structs.Paste, error) {
	return db.db.GetPagedPastes(page)
}

func (db *Database) GetPastesPages() int {
	return db.db.GetPastesPages()
}

// Initialize initializes connection to database.
func (db *Database) Initialize() {
	c.Logger.Info().Msg("Initializing database connection...")

	if c.Config.Database.Type == "mysql" {
		mysql.New(c)
	} else if c.Config.Database.Type == "flatfiles" {
		flatfiles.New(c)
	} else if c.Config.Database.Type == "postgresql" {
		postgresql.New(c)
	} else {
		c.Logger.Fatal().Str("type", c.Config.Database.Type).Msg("Unknown database type")
	}

	go db.cleanup()
}

func (db *Database) RegisterDialect(di dialectinterface.Interface) {
	db.db = di
	db.db.Initialize()
}

func (db *Database) SavePaste(p *structs.Paste) (int64, error) {
	return db.db.SavePaste(p)
}

func (db *Database) Shutdown() {
	db.db.Shutdown()
}
