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

package postgresql

import (
	// stdlib
	"database/sql"
	"fmt"
	"time"

	// local
	"gitlab.com/pztrn/fastpastebin/database/dialects/postgresql/migrations"
	"gitlab.com/pztrn/fastpastebin/pastes/model"

	// other
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database is a PostgreSQL connection controlling structure.
type Database struct {
	db *sqlx.DB
}

// Checks if queries can be run on known connection and reestablish it
// if not.
func (db *Database) check() {
	_, err := db.db.Exec("SELECT 1")
	if err != nil {
		db.Initialize()
	}
}

func (db *Database) GetDatabaseConnection() *sql.DB {
	db.check()
	return db.db.DB
}

// GetPaste returns a single paste by ID.
func (db *Database) GetPaste(pasteID int) (*pastesmodel.Paste, error) {
	db.check()
	p := &pastesmodel.Paste{}
	err := db.db.Get(p, db.db.Rebind("SELECT * FROM pastes WHERE id=$1"), pasteID)
	if err != nil {
		return nil, err
	}

	// We're aware of timezone in PostgreSQL, so SELECT will return
	// timestamps in server's local timezone. We should convert them.
	loc, _ := time.LoadLocation("UTC")

	utcCreatedAt := p.CreatedAt.In(loc)
	p.CreatedAt = &utcCreatedAt

	return p, nil
}

func (db *Database) GetPagedPastes(page int) ([]pastesmodel.Paste, error) {
	db.check()
	var pastesRaw []pastesmodel.Paste
	var pastes []pastesmodel.Paste

	// Pagination.
	var startPagination = 0
	if page > 1 {
		startPagination = (page - 1) * c.Config.Pastes.Pagination
	}

	err := db.db.Select(&pastesRaw, db.db.Rebind("SELECT * FROM pastes WHERE private != true ORDER BY id DESC LIMIT $1 OFFSET $2"), c.Config.Pastes.Pagination, startPagination)
	if err != nil {
		return nil, err
	}

	// We're aware of timezone in PostgreSQL, so SELECT will return
	// timestamps in server's local timezone. We should convert them.
	loc, _ := time.LoadLocation("UTC")

	for _, paste := range pastesRaw {
		if !paste.IsExpired() {
			utcCreatedAt := paste.CreatedAt.In(loc)
			paste.CreatedAt = &utcCreatedAt
			pastes = append(pastes, paste)
		}
	}

	return pastes, nil
}

func (db *Database) GetPastesPages() int {
	db.check()
	var pastesRaw []pastesmodel.Paste
	var pastes []pastesmodel.Paste
	err := db.db.Get(&pastesRaw, "SELECT * FROM pastes WHERE private != true")
	if err != nil {
		return 1
	}

	// Check if pastes isn't expired.
	for _, paste := range pastesRaw {
		if !paste.IsExpired() {
			pastes = append(pastes, paste)
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

	var userpass = ""
	if c.Config.Database.Password == "" {
		userpass = c.Config.Database.Username
	} else {
		userpass = c.Config.Database.Username + ":" + c.Config.Database.Password
	}

	dbConnString := fmt.Sprintf("postgres://%s@%s:%s/%s?connect_timeout=10&fallback_application_name=fastpastebin&sslmode=disable", userpass, c.Config.Database.Address, c.Config.Database.Port, c.Config.Database.Database)
	c.Logger.Debug().Msgf("Database connection string: %s", dbConnString)

	dbConn, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		c.Logger.Panic().Msgf("Failed to connect to database: %s", err.Error())
	}

	c.Logger.Info().Msg("Database connection established")
	db.db = dbConn

	// Perform migrations.
	migrations.New(c)
	migrations.Migrate()
}

func (db *Database) SavePaste(p *pastesmodel.Paste) (int64, error) {
	db.check()
	stmt, err := db.db.PrepareNamed("INSERT INTO pastes (title, data, created_at, keep_for, keep_for_unit_type, language, private, password, password_salt) VALUES (:title, :data, :created_at, :keep_for, :keep_for_unit_type, :language, :private, :password, :password_salt) RETURNING id")
	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.Get(&id, p)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) Shutdown() {
	err := db.db.Close()
	if err != nil {
		c.Logger.Error().Msgf("Failed to close database connection: %s", err.Error())
	}
}
