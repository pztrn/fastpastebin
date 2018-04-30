package http

import (
	// stdlib
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	// local
	"github.com/pztrn/fastpastebin/api/http/static"
	"github.com/pztrn/fastpastebin/models"

	// other
	"github.com/labstack/echo"
)

var (
	regexInts = regexp.MustCompile("[0-9]+")
)

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
	paste := &models.Paste{ID: pasteID}
	err1 := paste.GetByID(c.Database.GetDatabaseConnection())
	if err1 != nil {
		c.Logger.Error().Msgf("Failed to get paste #%d from database: %s", pasteID, err1.Error())
		return ec.HTML(http.StatusBadRequest, string(errhtml))
	}

	pasteHTML, err2 := static.ReadFile("paste.html")
	if err2 != nil {
		return ec.String(http.StatusNotFound, "parse.html wasn't found!")
	}

	pasteHTMLAsString := strings.Replace(string(pasteHTML), "{pastedata}", paste.Data, 1)

	return ec.HTML(http.StatusOK, string(pasteHTMLAsString))
}

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

	if !strings.ContainsAny(params["paste-keep-for"][0], "Mmhd") {
		c.Logger.Debug().Msgf("'Keep paste for' field have invalid value: %s", params["paste-keep-for"][0])
		return ec.HTML(http.StatusBadRequest, string(errhtml))
	}

	paste := &models.Paste{
		Title: params["paste-title"][0],
		Data:  params["paste-contents"][0],
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
		return ec.HTML(http.StatusBadRequest, string(errhtml))
	}
	paste.KeepFor = keepFor

	keepForUnitRaw := keepForUnitRegex.FindAllString(params["paste-keep-for"][0], 1)[0]
	keepForUnit := models.PASTE_KEEPS_CORELLATION[keepForUnitRaw]
	paste.KeepForUnitType = keepForUnit

	id, err2 := paste.Save(c.Database.GetDatabaseConnection())
	if err2 != nil {
		c.Logger.Debug().Msgf("Failed to save paste: %s", err2.Error())
		return ec.HTML(http.StatusBadRequest, string(errhtml))
	}

	newPasteIDAsString := strconv.FormatInt(id, 10)
	c.Logger.Debug().Msgf("Paste saved, URL: /paste/" + newPasteIDAsString)
	return ec.Redirect(http.StatusFound, "/paste/"+newPasteIDAsString)
}

func pastesGET(ec echo.Context) error {
	pasteListHTML, err1 := static.ReadFile("pastelist_list.html")
	if err1 != nil {
		return ec.String(http.StatusNotFound, "pastelist_list.html wasn't found!")
	}

	pasteElementHTML, err2 := static.ReadFile("pastelist_paste.html")
	if err2 != nil {
		return ec.String(http.StatusNotFound, "pastelist_paste.html wasn't found!")
	}

	pageFromParamRaw := ec.Param("id")
	var page = 1
	if pageFromParamRaw != "" {
		pageRaw := regexInts.FindAllString(pageFromParamRaw, 1)[0]
		page, _ = strconv.Atoi(pageRaw)
	}

	c.Logger.Debug().Msgf("Requested page #%d", page)

	p := &models.Paste{}
	// Get pastes IDs.
	pastes, err3 := p.GetPagedPastes(c.Database.GetDatabaseConnection(), page)
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

	return ec.HTML(http.StatusOK, string(pasteListHTMLAsString))
}
