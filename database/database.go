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
	"fmt"

	// other
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Database represents control structure for database connection.
type Database struct {
	db *sqlx.DB
}

// GetDatabaseConnection returns current database connection.
func (db *Database) GetDatabaseConnection() *sqlx.DB {
	return db.db
}

// Initialize initializes connection to database.
func (db *Database) Initialize() {
	c.Logger.Info().Msg("Initializing database connection...")

	// There might be only user, without password. MySQL/MariaDB driver
	// in DSN wants "user" or "user:password", "user:" is invalid.
	var userpass = ""
	if c.Config.Database.Password == "" {
		userpass = c.Config.Database.Username
	} else {
		userpass = c.Config.Database.Username + ":" + c.Config.Database.Password
	}

	dbConnString := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&collation=utf8mb4_unicode_ci&charset=utf8mb4", userpass, c.Config.Database.Address, c.Config.Database.Port, c.Config.Database.Database)
	c.Logger.Debug().Msgf("Database connection string: %s", dbConnString)

	dbConn, err := sqlx.Connect("mysql", dbConnString)
	if err != nil {
		c.Logger.Panic().Msgf("Failed to connect to database: %s", err.Error())
	}

	// Force UTC for current connection.
	_ = dbConn.MustExec("SET @@session.time_zone='+00:00';")

	c.Logger.Info().Msg("Database connection established")
	db.db = dbConn
}
