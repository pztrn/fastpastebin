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

package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/mysql/migrations"
	"go.dev.pztrn.name/fastpastebin/internal/structs"
)

// Database is a MySQL/MariaDB connection controlling structure.
type Database struct {
	db *sqlx.DB
}

// Checks if queries can be run on known connection and reestablish it
// if not.
func (db *Database) check() {
	if db.db != nil {
		_, err := db.db.Exec("SELECT 1")
		if err != nil {
			db.Initialize()
		}
	} else {
		db.Initialize()
	}
}

// DeletePaste deletes paste from database.
func (db *Database) DeletePaste(pasteID int) error {
	db.check()

	_, err := db.db.Exec(db.db.Rebind("DELETE FROM pastes WHERE id=?"), pasteID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetDatabaseConnection() *sql.DB {
	db.check()

	if db.db != nil {
		return db.db.DB
	}

	return nil
}

// GetPaste returns a single paste by ID.
func (db *Database) GetPaste(pasteID int) (*structs.Paste, error) {
	db.check()

	p := &structs.Paste{}

	err := db.db.Get(p, db.db.Rebind("SELECT * FROM `pastes` WHERE id=?"), pasteID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *Database) GetPagedPastes(page int) ([]structs.Paste, error) {
	db.check()

	var (
		pastesRaw []structs.Paste
		pastes    []structs.Paste
	)

	// Pagination.
	startPagination := 0
	if page > 1 {
		startPagination = (page - 1) * c.Config.Pastes.Pagination
	}

	err := db.db.Select(&pastesRaw, db.db.Rebind("SELECT * FROM `pastes` WHERE private != true ORDER BY id DESC LIMIT ? OFFSET ?"), c.Config.Pastes.Pagination, startPagination)
	if err != nil {
		return nil, err
	}

	for i := range pastesRaw {
		if !pastesRaw[i].IsExpired() {
			pastes = append(pastes, pastesRaw[i])
		}
	}

	return pastes, nil
}

func (db *Database) GetPastesPages() int {
	db.check()

	var (
		pastesRaw []structs.Paste
		pastes    []structs.Paste
	)

	err := db.db.Get(&pastesRaw, "SELECT * FROM `pastes` WHERE private != true")
	if err != nil {
		return 1
	}

	// Check if pastes isn't expired.
	for i := range pastesRaw {
		if !pastesRaw[i].IsExpired() {
			pastes = append(pastes, pastesRaw[i])
		}
	}

	// Calculate pages.
	pages := len(pastes) / c.Config.Pastes.Pagination
	// Check if we have any remainder. Add 1 to pages count if so.
	if len(pastes)%c.Config.Pastes.Pagination > 0 {
		pages++
	}

	return pages
}

// Initialize initializes MySQL/MariaDB connection.
func (db *Database) Initialize() {
	c.Logger.Info().Msg("Initializing database connection...")

	// There might be only user, without password. MySQL/MariaDB driver
	// in DSN wants "user" or "user:password", "user:" is invalid.
	var userpass string
	if c.Config.Database.Password == "" {
		userpass = c.Config.Database.Username
	} else {
		userpass = c.Config.Database.Username + ":" + c.Config.Database.Password
	}

	dbConnString := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&collation=utf8mb4_unicode_ci&charset=utf8mb4", userpass, c.Config.Database.Address, c.Config.Database.Port, c.Config.Database.Database)
	c.Logger.Debug().Str("DSN", dbConnString).Msgf("Database connection string composed")

	dbConn, err := sqlx.Connect("mysql", dbConnString)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to connect to database")

		return
	}

	// Force UTC for current connection.
	_ = dbConn.MustExec("SET @@session.time_zone='+00:00';")

	c.Logger.Info().Msg("Database connection established")

	db.db = dbConn

	// Perform migrations.
	migrations.New(c)
	migrations.Migrate()
}

func (db *Database) SavePaste(p *structs.Paste) (int64, error) {
	db.check()

	result, err := db.db.NamedExec("INSERT INTO `pastes` (title, data, created_at, keep_for, keep_for_unit_type, language, private, password, password_salt) VALUES (:title, :data, :created_at, :keep_for, :keep_for_unit_type, :language, :private, :password, :password_salt)", p)
	if err != nil {
		return 0, err
	}

	ID, err1 := result.LastInsertId()
	if err1 != nil {
		return 0, err
	}

	return ID, nil
}

func (db *Database) Shutdown() {
	if db.db != nil {
		err := db.db.Close()
		if err != nil {
			c.Logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}
}
