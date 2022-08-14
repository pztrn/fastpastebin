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
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/flatfiles"
	"go.dev.pztrn.name/fastpastebin/internal/pagination"
	"go.dev.pztrn.name/fastpastebin/internal/templater"
)

// GET for "/pastes/", a list of publicly available pastes.
// Web interface version.
func pastesGET(ectx echo.Context) error {
	// We should check if database connection available.
	dbConn := ctx.Database.GetDatabaseConnection()
	if ctx.Config.Database.Type != flatfiles.FlatFileDialect && dbConn == nil {
		//nolint:wrapcheck
		return ectx.Redirect(http.StatusFound, "/database_not_available")
	}

	pageFromParamRaw := ectx.Param("page")

	page := 1

	if pageFromParamRaw != "" {
		pageRaw := regexInts.FindAllString(pageFromParamRaw, 1)[0]
		page, _ = strconv.Atoi(pageRaw)
	}

	ctx.Logger.Debug().Int("page", page).Msg("Requested page")

	// Get pastes IDs.
	pastes, err3 := ctx.Database.GetPagedPastes(page)
	ctx.Logger.Debug().Int("count", len(pastes)).Msg("Got pastes")

	pastesString := "No pastes to show."

	// Show "No pastes to show" on any error for now.
	if err3 != nil {
		ctx.Logger.Error().Err(err3).Msg("Failed to get pastes list from database")

		noPastesToShowTpl := templater.GetErrorTemplate(ectx, "No pastes to show.")

		//nolint:wrapcheck
		return ectx.HTML(http.StatusOK, noPastesToShowTpl)
	}

	if len(pastes) > 0 {
		pastesString = ""

		for _, paste := range pastes {
			pasteDataMap := make(map[string]string)
			pasteDataMap["pasteID"] = strconv.Itoa(paste.ID)
			pasteDataMap["pasteTitle"] = paste.Title
			pasteDataMap["pasteDate"] = paste.CreatedAt.Format("2006-01-02 @ 15:04:05") + " UTC"

			// Get max 4 lines of each paste.
			pasteDataSplitted := strings.Split(paste.Data, "\n")

			var pasteData string

			if len(pasteDataSplitted) < 4 {
				pasteData = paste.Data
			} else {
				pasteData = strings.Join(pasteDataSplitted[0:4], "\n")
			}

			pasteDataMap["pasteData"] = pasteData
			pasteTpl := templater.GetRawTemplate(ectx, "pastelist_paste.html", pasteDataMap)

			pastesString += pasteTpl
		}
	}

	// Pagination.
	pages := ctx.Database.GetPastesPages()
	ctx.Logger.Debug().Int("total pages", pages).Int("current page", page).Msg("Paging data")
	paginationHTML := pagination.CreateHTML(page, pages, "/pastes/")

	pasteListTpl := templater.GetTemplate(ectx, "pastelist_list.html", map[string]string{"pastes": pastesString, "pagination": paginationHTML})

	//nolint:wrapcheck
	return ectx.HTML(http.StatusOK, pasteListTpl)
}
