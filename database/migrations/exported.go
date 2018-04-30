package migrations

import (
	// local
	"github.com/pztrn/fastpastebin/context"

	// other
	//"github.com/jmoiron/sqlx"
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

	goose.SetDialect("mysql")
	goose.AddNamedMigration("1_initial.go", InitialUp, nil)
	// Add new migrations BEFORE this message.

	dbConn := c.Database.GetDatabaseConnection()
	err := goose.Up(dbConn.DB, ".")
	if err != nil {
		c.Logger.Panic().Msgf("Failed to migrate database to latest version: %s", err.Error())
	}
}
