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
	"regexp"
	"strconv"
	"strings"
	"time"

	// local
	"github.com/pztrn/fastpastebin/api/http/static"
	"github.com/pztrn/fastpastebin/pagination"

	// other
	"github.com/alecthomas/chroma/lexers"
	"github.com/labstack/echo"
)

var (
	regexInts = regexp.MustCompile("[0-9]+")
)

// GET for "/paste/PASTE_ID".
func pasteGET(ec echo.Context) error {
	errhtml, err := static.ReadFile("error.html")
	if err != nil {
		return ec.String(http.StatusNotFound, "error.html wasn't found!")
	}

	pasteIDRaw := ec.Param("id")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	c.Logger.Debug().Msgf("Requesting paste #%+v", pasteID)

	// Get paste.
	paste, err1 := GetByID(pasteID)
	if err1 != nil {
		c.Logger.Error().Msgf("Failed to get paste #%d from database: %s", pasteID, err1.Error())
		errhtmlAsString := strings.Replace(string(errhtml), "{error}", "Paste #"+strconv.Itoa(pasteID)+" not found", 1)
		return ec.HTML(http.StatusBadRequest, errhtmlAsString)
	}

	pasteHTML, err2 := static.ReadFile("paste.html")
	if err2 != nil {
		return ec.String(http.StatusNotFound, "parse.html wasn't found!")
	}

	pasteHTMLAsString := strings.Replace(string(pasteHTML), "{pastedata}", paste.Data, 1)

	return ec.HTML(http.StatusOK, string(pasteHTMLAsString))
}

// POST for "/paste/" which will create new paste and redirect to
// "/pastes/CREATED_PASTE_ID".
func pastePOST(ec echo.Context) error {
	errhtml, err := static.ReadFile("error.html")
	if err != nil {
		return ec.String(http.StatusNotFound, "error.html wasn't found!")
	}

	params, err := ec.FormParams()
	if err != nil {
		c.Logger.Debug().Msg("No form parameters passed")
		return ec.HTML(http.StatusBadRequest, string(errhtml))
	}
	c.Logger.Debug().Msgf("Received parameters: %+v", params)

	// Do nothing if paste contents is empty.
	if len(params["paste-contents"][0]) == 0 {
		c.Logger.Debug().Msg("Empty paste submitted, ignoring")
		errhtmlAsString := strings.Replace(string(errhtml), "{error}", "Empty pastes aren't allowed.", 1)
		return ec.HTML(http.StatusBadRequest, errhtmlAsString)
	}

	if !strings.ContainsAny(params["paste-keep-for"][0], "Mmhd") {
		c.Logger.Debug().Msgf("'Keep paste for' field have invalid value: %s", params["paste-keep-for"][0])
		errhtmlAsString := strings.Replace(string(errhtml), "{error}", "Invalid 'Paste should be available for' parameter passed. Please do not try to hack us ;).", 1)
		return ec.HTML(http.StatusBadRequest, errhtmlAsString)
	}

	paste := &Paste{
		Title:    params["paste-title"][0],
		Data:     params["paste-contents"][0],
		Language: params["paste-language"][0],
	}

	// Paste creation time in UTC.
	createdAt := time.Now().UTC()
	paste.CreatedAt = &createdAt

	// Parse "keep for" field.

	// Get integers and strings separately.
	keepForUnitRegex := regexp.MustCompile("[Mmhd]")

	keepForRaw := regexInts.FindAllString(params["paste-keep-for"][0], 1)[0]
	keepFor, err1 := strconv.Atoi(keepForRaw)
	if err1 != nil {
		c.Logger.Debug().Msgf("Failed to parse 'Keep for' integer: %s", err1.Error())
		errhtmlAsString := strings.Replace(string(errhtml), "{error}", "Invalid 'Paste should be available for' parameter passed. Please do not try to hack us ;).", 1)
		return ec.HTML(http.StatusBadRequest, errhtmlAsString)
	}
	paste.KeepFor = keepFor

	keepForUnitRaw := keepForUnitRegex.FindAllString(params["paste-keep-for"][0], 1)[0]
	keepForUnit := PASTE_KEEPS_CORELLATION[keepForUnitRaw]
	paste.KeepForUnitType = keepForUnit

	// Try to autodetect if it was selected.
	if params["paste-language"][0] == "autodetect" {
		lexer := lexers.Analyse(params["paste-language"][0])
		if lexer != nil {
			paste.Language = lexer.Config().Name
		} else {
			paste.Language = "text"
		}
	}

	id, err2 := Save(paste)
	if err2 != nil {
		c.Logger.Debug().Msgf("Failed to save paste: %s", err2.Error())
		errhtmlAsString := strings.Replace(string(errhtml), "{error}", "Failed to save paste. Please, try again later.", 1)
		return ec.HTML(http.StatusBadRequest, errhtmlAsString)
	}

	newPasteIDAsString := strconv.FormatInt(id, 10)
	c.Logger.Debug().Msgf("Paste saved, URL: /paste/" + newPasteIDAsString)
	return ec.Redirect(http.StatusFound, "/paste/"+newPasteIDAsString)
}

// GET for "/pastes/", a list of publicly available pastes.
func pastesGET(ec echo.Context) error {
	pasteListHTML, err1 := static.ReadFile("pastelist_list.html")
	if err1 != nil {
		return ec.String(http.StatusNotFound, "pastelist_list.html wasn't found!")
	}

	pasteElementHTML, err2 := static.ReadFile("pastelist_paste.html")
	if err2 != nil {
		return ec.String(http.StatusNotFound, "pastelist_paste.html wasn't found!")
	}

	pageFromParamRaw := ec.Param("page")
	var page = 1
	if pageFromParamRaw != "" {
		pageRaw := regexInts.FindAllString(pageFromParamRaw, 1)[0]
		page, _ = strconv.Atoi(pageRaw)
	}

	c.Logger.Debug().Msgf("Requested page #%d", page)

	// Get pastes IDs.
	pastes, err3 := GetPagedPastes(page)
	c.Logger.Debug().Msgf("Got %d pastes", len(pastes))

	var pastesString = "No pastes to show."

	// Show "No pastes to show" on any error for now.
	if err3 != nil {
		c.Logger.Error().Msgf("Failed to get pastes list from database: %s", err3.Error())
		pasteListHTMLAsString := strings.Replace(string(pasteListHTML), "{pastes}", pastesString, 1)
		return ec.HTML(http.StatusOK, string(pasteListHTMLAsString))
	}

	if len(pastes) > 0 {
		pastesString = ""
		for i := range pastes {
			pasteString := strings.Replace(string(pasteElementHTML), "{pasteID}", strconv.Itoa(pastes[i].ID), 2)
			pasteString = strings.Replace(pasteString, "{pasteTitle}", pastes[i].Title, 1)
			pasteString = strings.Replace(pasteString, "{pasteDate}", pastes[i].CreatedAt.Format("2006-01-02 @ 15:04:05"), 1)

			// Get max 4 lines of each paste.
			pasteDataSplitted := strings.Split(pastes[i].Data, "\n")
			var pasteData = ""
			if len(pasteDataSplitted) < 4 {
				pasteData = pastes[i].Data
			} else {
				pasteData = strings.Join(pasteDataSplitted[0:4], "\n")
			}
			pasteString = strings.Replace(pasteString, "{pasteData}", pasteData, 1)

			pastesString += pasteString
		}
	}

	pasteListHTMLAsString := strings.Replace(string(pasteListHTML), "{pastes}", pastesString, 1)

	pages := GetPastesPages()
	c.Logger.Debug().Msgf("Total pages: %d, current: %d", pages, page)
	paginationHTML := pagination.CreateHTML(page, pages, "/pastes/")
	pasteListHTMLAsString = strings.Replace(pasteListHTMLAsString, "{pagination}", paginationHTML, -1)

	return ec.HTML(http.StatusOK, string(pasteListHTMLAsString))
}
