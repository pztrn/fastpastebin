package databaseinterface

import (
	// other
	"github.com/jmoiron/sqlx"
)

type Interface interface {
	GetDatabaseConnection() *sqlx.DB
	Initialize()
}
