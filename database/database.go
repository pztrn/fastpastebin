package database

import (
	// stdlib
	"fmt"

	// other
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func (db *Database) GetDatabaseConnection() *sqlx.DB {
	return db.db
}

func (db *Database) Initialize() {
	c.Logger.Info().Msg("Initializing database connection...")

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
