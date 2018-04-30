package http

import (
	// stdlib
	"net/http"

	// local
	"github.com/pztrn/fastpastebin/api/http/static"

	// other
	"github.com/labstack/echo"
)

func indexGet(ec echo.Context) error {
	html, err := static.ReadFile("index.html")
	if err != nil {
		return ec.String(http.StatusNotFound, "index.html wasn't found!")
	}

	return ec.HTML(http.StatusOK, string(html))
}
