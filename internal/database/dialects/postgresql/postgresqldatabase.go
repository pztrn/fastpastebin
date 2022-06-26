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

// nolint:gci
import (
	"database/sql"
	"fmt"
	"net"
	"time"

	// PostgreSQL driver.
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/postgresql/migrations"
	"go.dev.pztrn.name/fastpastebin/internal/structs"
)

// Database is a PostgreSQL connection controlling structure.
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
		// nolint:wrapcheck
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

	// nolint:exhaustruct
	paste := &structs.Paste{}

	err := db.db.Get(paste, db.db.Rebind("SELECT * FROM pastes WHERE id=$1"), pasteID)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	// We're aware of timezone in PostgreSQL, so SELECT will return
	// timestamps in server's local timezone. We should convert them.
	loc, _ := time.LoadLocation("UTC")

	utcCreatedAt := paste.CreatedAt.In(loc)
	paste.CreatedAt = &utcCreatedAt

	return paste, nil
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
		startPagination = (page - 1) * ctx.Config.Pastes.Pagination
	}

	err := db.db.Select(&pastesRaw, db.db.Rebind("SELECT * FROM pastes WHERE private != true ORDER BY id DESC LIMIT $1 OFFSET $2"), ctx.Config.Pastes.Pagination, startPagination)
	if err != nil {
		// nolint:wrapcheck
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

	var (
		pastesRaw []structs.Paste
		pastes    []structs.Paste
	)

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
	pages := len(pastes) / ctx.Config.Pastes.Pagination
	// Check if we have any remainder. Add 1 to pages count if so.
	if len(pastes)%ctx.Config.Pastes.Pagination > 0 {
		pages++
	}

	return pages
}

// Initialize initializes MySQL/MariaDB connection.
func (db *Database) Initialize() {
	ctx.Logger.Info().Msg("Initializing database connection...")

	var userpass string
	if ctx.Config.Database.Password == "" {
		userpass = ctx.Config.Database.Username
	} else {
		userpass = ctx.Config.Database.Username + ":" + ctx.Config.Database.Password
	}

	dbConnString := fmt.Sprintf("postgres://%s@%s/%s?connect_timeout=10&fallback_application_name=fastpastebin&sslmode=disable", userpass, net.JoinHostPort(ctx.Config.Database.Address, ctx.Config.Database.Port), ctx.Config.Database.Database)
	ctx.Logger.Debug().Str("DSN", dbConnString).Msg("Database connection string composed")

	dbConn, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		ctx.Logger.Error().Err(err).Msg("Failed to connect to database")

		return
	}

	ctx.Logger.Info().Msg("Database connection established")

	db.db = dbConn

	// Perform migrations.
	migrations.New(ctx)
	migrations.Migrate()
}

func (db *Database) SavePaste(paste *structs.Paste) (int64, error) {
	db.check()

	stmt, err := db.db.PrepareNamed("INSERT INTO pastes (title, data, created_at, keep_for, keep_for_unit_type, language, private, password, password_salt) VALUES (:title, :data, :created_at, :keep_for, :keep_for_unit_type, :language, :private, :password, :password_salt) RETURNING id")
	if err != nil {
		// nolint:wrapcheck
		return 0, err
	}

	var newPasteID int64

	err = stmt.Get(&newPasteID, paste)
	if err != nil {
		// nolint:wrapcheck
		return 0, err
	}

	return newPasteID, nil
}

func (db *Database) Shutdown() {
	if db.db != nil {
		err := db.db.Close()
		if err != nil {
			ctx.Logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}
}
