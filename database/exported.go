package database

import (
	// local
	"github.com/pztrn/fastpastebin/context"
	"github.com/pztrn/fastpastebin/database/interface"
)

var (
	c *context.Context
	d *Database
)

func New(cc *context.Context) {
	c = cc
	d = &Database{}
	c.RegisterDatabaseInterface(databaseinterface.Interface(Handler{}))
}
