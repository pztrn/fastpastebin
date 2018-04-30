package database

import (
	// other
	"github.com/jmoiron/sqlx"
)

type Handler struct{}

func (dbh Handler) GetDatabaseConnection() *sqlx.DB {
	return d.GetDatabaseConnection()
}

func (dbh Handler) Initialize() {
	d.Initialize()
}
