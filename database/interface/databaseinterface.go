package databaseinterface

import (
	// other
	"github.com/jmoiron/sqlx"
)

// Interface represents database interface which is available to all
// parts of application and registers with context.Context.
type Interface interface {
	GetDatabaseConnection() *sqlx.DB
	Initialize()
}
