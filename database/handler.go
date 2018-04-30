package database

import (
	// other
	"github.com/jmoiron/sqlx"
)

// Handler is an interfaceable structure that proxifies calls from anyone
// to Database structure.
type Handler struct{}

// GetDatabaseConnection returns current database connection.
func (dbh Handler) GetDatabaseConnection() *sqlx.DB {
	return d.GetDatabaseConnection()
}

// Initialize initializes connection to database.
func (dbh Handler) Initialize() {
	d.Initialize()
}
