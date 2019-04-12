// Fast Paste Bin - uberfast and easy-to-use pastebin.
//
// Copyright (c) 2018, Stanislav N. aka pztrn and Fast Paste Bin
// developers.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject
// to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package pastes

import (
	// stdlib
	"net/http"
	"strconv"
	"strings"

	// local
	"gitlab.com/pztrn/fastpastebin/internal/pagination"
	"gitlab.com/pztrn/fastpastebin/internal/templater"

	// other
	"github.com/labstack/echo"
)

// GET for "/pastes/", a list of publicly available pastes.
// Web inteface version.
func pastesGET(ec echo.Context) error {
	// We should check if database connection available.
	dbConn := c.Database.GetDatabaseConnection()
	if c.Config.Database.Type != "flatfiles" && dbConn == nil {
		return ec.Redirect(http.StatusFound, "/database_not_available")
	}

	pageFromParamRaw := ec.Param("page")
	var page = 1
	if pageFromParamRaw != "" {
		pageRaw := regexInts.FindAllString(pageFromParamRaw, 1)[0]
		page, _ = strconv.Atoi(pageRaw)
	}

	c.Logger.Debug().Msgf("Requested page #%d", page)

	// Get pastes IDs.
	pastes, err3 := c.Database.GetPagedPastes(page)
	c.Logger.Debug().Msgf("Got %d pastes", len(pastes))

	var pastesString = "No pastes to show."

	// Show "No pastes to show" on any error for now.
	if err3 != nil {
		c.Logger.Error().Msgf("Failed to get pastes list from database: %s", err3.Error())
		noPastesToShowTpl := templater.GetErrorTemplate(ec, "No pastes to show.")
		return ec.HTML(http.StatusOK, noPastesToShowTpl)
	}

	if len(pastes) > 0 {
		pastesString = ""
		for i := range pastes {
			pasteDataMap := make(map[string]string)
			pasteDataMap["pasteID"] = strconv.Itoa(pastes[i].ID)
			pasteDataMap["pasteTitle"] = pastes[i].Title
			pasteDataMap["pasteDate"] = pastes[i].CreatedAt.Format("2006-01-02 @ 15:04:05") + " UTC"

			// Get max 4 lines of each paste.
			pasteDataSplitted := strings.Split(pastes[i].Data, "\n")
			var pasteData = ""
			if len(pasteDataSplitted) < 4 {
				pasteData = pastes[i].Data
			} else {
				pasteData = strings.Join(pasteDataSplitted[0:4], "\n")
			}

			pasteDataMap["pasteData"] = pasteData
			pasteTpl := templater.GetRawTemplate(ec, "pastelist_paste.html", pasteDataMap)

			pastesString += pasteTpl
		}
	}

	// Pagination.
	pages := c.Database.GetPastesPages()
	c.Logger.Debug().Msgf("Total pages: %d, current: %d", pages, page)
	paginationHTML := pagination.CreateHTML(page, pages, "/pastes/")

	pasteListTpl := templater.GetTemplate(ec, "pastelist_list.html", map[string]string{"pastes": pastesString, "pagination": paginationHTML})

	return ec.HTML(http.StatusOK, string(pasteListTpl))
}
