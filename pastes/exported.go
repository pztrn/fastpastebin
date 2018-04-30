package pastes

import (
	// local
	"github.com/pztrn/fastpastebin/context"
)

var (
	c *context.Context
)

// New initializes pastes package and adds neccessary HTTP and API
// endpoints.
func New(cc *context.Context) {
	c = cc

	// New paste.
	c.Echo.POST("/paste/", pastePOST)

	// Show paste.
	c.Echo.GET("/paste/:id", pasteGET)

	// Pastes list.
	c.Echo.GET("/pastes/", pastesGET)
	c.Echo.GET("/pastes/:page", pastesGET)
}
