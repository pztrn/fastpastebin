package json

import (
	// local
	"github.com/pztrn/fastpastebin/context"
)

var (
	c *context.Context
)

// New initializes basic JSON API.
func New(cc *context.Context) {
	c = cc
	c.Logger.Info().Msg("Initializing JSON API...")
}
