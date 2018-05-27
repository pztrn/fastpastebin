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
	"bytes"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	// local
	"github.com/pztrn/fastpastebin/captcha"
	"github.com/pztrn/fastpastebin/pagination"
	"github.com/pztrn/fastpastebin/pastes/model"
	"github.com/pztrn/fastpastebin/templater"

	// other
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	htmlfmt "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	//"github.com/dchest/captcha"
	"github.com/labstack/echo"
)

var (
	regexInts = regexp.MustCompile("[0-9]+")
)

// GET for "/paste/PASTE_ID" and "/paste/PASTE_ID/TIMESTAMP" (private pastes).
func pasteGET(ec echo.Context) error {
	pasteIDRaw := ec.Param("id")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	c.Logger.Debug().Msgf("Requesting paste #%+v", pasteID)

	// Get paste.
	paste, err1 := c.Database.GetPaste(pasteID)
	if err1 != nil {
		c.Logger.Error().Msgf("Failed to get paste #%d: %s", pasteID, err1.Error())
		errtpl := templater.GetErrorTemplate(ec, "Paste #"+strconv.Itoa(pasteID)+" not found")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	if paste.IsExpired() {
		c.Logger.Error().Msgf("Paste #%d is expired", pasteID)
		errtpl := templater.GetErrorTemplate(ec, "Paste #"+strconv.Itoa(pasteID)+" not found")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	// Check if we have a private paste and it's parameters are correct.
	if paste.Private {
		tsProvidedStr := ec.Param("timestamp")
		tsProvided, err2 := strconv.ParseInt(tsProvidedStr, 10, 64)
		if err2 != nil {
			c.Logger.Error().Msgf("Invalid timestamp '%s' provided for getting private paste #%d: %s", tsProvidedStr, pasteID, err2.Error())
			errtpl := templater.GetErrorTemplate(ec, "Paste #"+strconv.Itoa(pasteID)+" not found")
			return ec.HTML(http.StatusBadRequest, errtpl)
		}
		pasteTs := paste.CreatedAt.Unix()
		if tsProvided != pasteTs {
			c.Logger.Error().Msgf("Incorrect timestamp '%v' provided for private paste #%d, waiting for %v", tsProvidedStr, pasteID, strconv.FormatInt(pasteTs, 10))
			errtpl := templater.GetErrorTemplate(ec, "Paste #"+strconv.Itoa(pasteID)+" not found")
			return ec.HTML(http.StatusBadRequest, errtpl)
		}
	}

	if paste.Private && paste.Password != "" {
		// Check if cookie for this paste is defined. This means that user
		// previously successfully entered a password.
		cookie, err := ec.Cookie("PASTE-" + strconv.Itoa(pasteID))
		if err != nil {
			// No cookie, redirect to auth page.
			c.Logger.Info().Msg("Tried to access passworded paste without autorization, redirecting to auth page...")
			return ec.Redirect(http.StatusMovedPermanently, "/paste/"+pasteIDRaw+"/"+ec.Param("timestamp")+"/verify")
		}

		// Generate cookie value to check.
		cookieValue := paste.GenerateCryptedCookieValue()

		if cookieValue != cookie.Value {
			c.Logger.Info().Msg("Invalid cookie, redirecting to auth page...")
			return ec.Redirect(http.StatusMovedPermanently, "/paste/"+pasteIDRaw+"/"+ec.Param("timestamp")+"/verify")
		}

		// If all okay - do nothing :)
	}

	// Format paste data map.
	pasteData := make(map[string]string)
	pasteData["pasteTitle"] = paste.Title
	pasteData["pasteID"] = strconv.Itoa(paste.ID)
	pasteData["pasteDate"] = paste.CreatedAt.Format("2006-01-02 @ 15:04:05") + " UTC"
	pasteData["pasteExpiration"] = paste.GetExpirationTime().Format("2006-01-02 @ 15:04:05") + " UTC"
	pasteData["pasteLanguage"] = paste.Language

	if paste.Private {
		pasteData["pasteType"] = "<span class='has-text-danger'>Private</span>"
		pasteData["pasteTs"] = strconv.FormatInt(paste.CreatedAt.Unix(), 10) + "/"
	} else {
		pasteData["pasteType"] = "<span class='has-text-success'>Public</span>"
		pasteData["pasteTs"] = ""
	}

	// Highlight.
	// Get lexer.
	lexer := lexers.Get(paste.Language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	// Tokenize paste data.
	lexered, err3 := lexer.Tokenise(nil, paste.Data)
	if err3 != nil {
		c.Logger.Error().Msgf("Failed to tokenize paste data: %s", err3.Error())
	}
	// Get style for HTML output.
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	// Get HTML formatter.
	formatter := chroma.Formatter(htmlfmt.New(htmlfmt.WithLineNumbers(), htmlfmt.LineNumbersInTable()))
	if formatter == nil {
		formatter = formatters.Fallback
	}
	// Create buffer and format into it.
	buf := new(bytes.Buffer)
	err4 := formatter.Format(buf, style, lexered)
	if err4 != nil {
		c.Logger.Error().Msgf("Failed to format paste data: %s", err4.Error())
	}
	pasteData["pastedata"] = buf.String()

	// Get template and format it.
	pasteHTML := templater.GetTemplate(ec, "paste.html", pasteData)

	return ec.HTML(http.StatusOK, pasteHTML)
}

// GET for "/paste/PASTE_ID/TIMESTAMP/verify" - a password verify page.
func pastePasswordedVerifyGet(ec echo.Context) error {
	pasteIDRaw := ec.Param("id")
	timestampRaw := ec.Param("timestamp")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])

	// Get paste.
	paste, err1 := c.Database.GetPaste(pasteID)
	if err1 != nil {
		c.Logger.Error().Msgf("Failed to get paste #%d: %s", pasteID, err1.Error())
		errtpl := templater.GetErrorTemplate(ec, "Paste #"+pasteIDRaw+" not found")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	// Check for auth cookie. If present - redirect to paste.
	cookie, err := ec.Cookie("PASTE-" + strconv.Itoa(pasteID))
	if err == nil {
		// No cookie, redirect to auth page.
		c.Logger.Debug().Msg("Paste cookie found, checking it...")

		// Generate cookie value to check.
		cookieValue := paste.GenerateCryptedCookieValue()

		if cookieValue == cookie.Value {
			c.Logger.Info().Msg("Valid cookie, redirecting to paste page...")
			return ec.Redirect(http.StatusMovedPermanently, "/paste/"+pasteIDRaw+"/"+ec.Param("timestamp"))
		}

		c.Logger.Debug().Msg("Invalid cookie, showing auth page")
	}

	// HTML data.
	htmlData := make(map[string]string)
	htmlData["pasteID"] = strconv.Itoa(pasteID)
	htmlData["pasteTimestamp"] = timestampRaw

	verifyHTML := templater.GetTemplate(ec, "passworded_paste_verify.html", htmlData)

	return ec.HTML(http.StatusOK, verifyHTML)
}

// POST for "/paste/PASTE_ID/TIMESTAMP/verify" - a password verify page.
func pastePasswordedVerifyPost(ec echo.Context) error {
	pasteIDRaw := ec.Param("id")
	timestampRaw := ec.Param("timestamp")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	c.Logger.Debug().Msgf("Requesting paste #%+v", pasteID)

	// Get paste.
	paste, err1 := c.Database.GetPaste(pasteID)
	if err1 != nil {
		c.Logger.Error().Msgf("Failed to get paste #%d: %s", pasteID, err1.Error())
		errtpl := templater.GetErrorTemplate(ec, "Paste #"+strconv.Itoa(pasteID)+" not found")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	params, err2 := ec.FormParams()
	if err2 != nil {
		c.Logger.Debug().Msg("No form parameters passed")
		errtpl := templater.GetErrorTemplate(ec, "Paste #"+strconv.Itoa(pasteID)+" not found")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	if paste.VerifyPassword(params["paste-password"][0]) {
		// Set cookie that this paste's password is verified and paste
		// can be viewed.
		cookie := new(http.Cookie)
		cookie.Name = "PASTE-" + strconv.Itoa(pasteID)
		cookie.Value = paste.GenerateCryptedCookieValue()
		cookie.Expires = time.Now().Add(24 * time.Hour)
		ec.SetCookie(cookie)

		return ec.Redirect(http.StatusFound, "/paste/"+strconv.Itoa(pasteID)+"/"+timestampRaw)
	}

	errtpl := templater.GetErrorTemplate(ec, "Invalid password. Please, try again.")
	return ec.HTML(http.StatusBadRequest, string(errtpl))
}

// POST for "/paste/" which will create new paste and redirect to
// "/pastes/CREATED_PASTE_ID".
func pastePOST(ec echo.Context) error {
	params, err := ec.FormParams()
	if err != nil {
		errtpl := templater.GetErrorTemplate(ec, "Cannot create empty paste")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}
	c.Logger.Debug().Msgf("Received parameters: %+v", params)

	// Do nothing if paste contents is empty.
	if len(params["paste-contents"][0]) == 0 {
		c.Logger.Debug().Msg("Empty paste submitted, ignoring")
		errtpl := templater.GetErrorTemplate(ec, "Empty pastes aren't allowed.")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	if !strings.ContainsAny(params["paste-keep-for"][0], "Mmhd") {
		c.Logger.Debug().Msgf("'Keep paste for' field have invalid value: %s", params["paste-keep-for"][0])
		errtpl := templater.GetErrorTemplate(ec, "Invalid 'Paste should be available for' parameter passed. Please do not try to hack us ;).")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	// Verify captcha.
	if !captcha.Verify(params["paste-captcha-id"][0], params["paste-captcha-solution"][0]) {
		c.Logger.Debug().Msgf("Invalid captcha solution for captcha ID '%s': %s", params["paste-captcha-id"][0], params["paste-captcha-solution"][0])
		errtpl := templater.GetErrorTemplate(ec, "Invalid captcha solution.")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	paste := &pastesmodel.Paste{
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
		errtpl := templater.GetErrorTemplate(ec, "Invalid 'Paste should be available for' parameter passed. Please do not try to hack us ;).")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}
	paste.KeepFor = keepFor

	keepForUnitRaw := keepForUnitRegex.FindAllString(params["paste-keep-for"][0], 1)[0]
	keepForUnit := pastesmodel.PASTE_KEEPS_CORELLATION[keepForUnitRaw]
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

	// Private paste?
	paste.Private = false
	privateCheckbox, privateCheckboxFound := params["paste-private"]
	pastePassword, pastePasswordFound := params["paste-password"]
	if privateCheckboxFound && privateCheckbox[0] == "on" || pastePasswordFound && pastePassword[0] != "" {
		paste.Private = true
	}

	if pastePassword[0] != "" {
		paste.CreatePassword(pastePassword[0])
	}

	id, err2 := c.Database.SavePaste(paste)
	if err2 != nil {
		c.Logger.Debug().Msgf("Failed to save paste: %s", err2.Error())
		errtpl := templater.GetErrorTemplate(ec, "Failed to save paste. Please, try again later.")
		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	newPasteIDAsString := strconv.FormatInt(id, 10)
	c.Logger.Debug().Msgf("Paste saved, URL: /paste/" + newPasteIDAsString)

	// Private pastes have it's timestamp in URL.
	if paste.Private {
		return ec.Redirect(http.StatusFound, "/paste/"+newPasteIDAsString+"/"+strconv.FormatInt(paste.CreatedAt.Unix(), 10))
	}

	return ec.Redirect(http.StatusFound, "/paste/"+newPasteIDAsString)
}

// GET for "/pastes/:id/raw", raw paste output.
func pasteRawGET(ec echo.Context) error {
	pasteIDRaw := ec.Param("id")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	c.Logger.Debug().Msgf("Requesting paste #%+v", pasteID)

	// Get paste.
	paste, err1 := c.Database.GetPaste(pasteID)
	if err1 != nil {
		c.Logger.Error().Msgf("Failed to get paste #%d from database: %s", pasteID, err1.Error())
		return ec.HTML(http.StatusBadRequest, "Paste #"+pasteIDRaw+" does not exist.")
	}

	if paste.IsExpired() {
		c.Logger.Error().Msgf("Paste #%d is expired", pasteID)
		return ec.HTML(http.StatusBadRequest, "Paste #"+pasteIDRaw+" does not exist.")
	}

	// Check if we have a private paste and it's parameters are correct.
	if paste.Private {
		tsProvidedStr := ec.Param("timestamp")
		tsProvided, err2 := strconv.ParseInt(tsProvidedStr, 10, 64)
		if err2 != nil {
			c.Logger.Error().Msgf("Invalid timestamp '%s' provided for getting private paste #%d: %s", tsProvidedStr, pasteID, err2.Error())
			return ec.String(http.StatusBadRequest, "Paste #"+pasteIDRaw+" not found")
		}
		pasteTs := paste.CreatedAt.Unix()
		if tsProvided != pasteTs {
			c.Logger.Error().Msgf("Incorrect timestamp '%v' provided for private paste #%d, waiting for %v", tsProvidedStr, pasteID, strconv.FormatInt(pasteTs, 10))
			return ec.String(http.StatusBadRequest, "Paste #"+pasteIDRaw+" not found")
		}
	}

	// ToDo: figure out how to handle passworded pastes here.
	// Return error for now.
	if paste.Password != "" {
		c.Logger.Error().Msgf("Cannot render paste #%d as raw: passworded paste. Patches welcome!", pasteID)
		return ec.String(http.StatusBadRequest, "Paste #"+pasteIDRaw+" not found")
	}

	return ec.String(http.StatusOK, paste.Data)
}

// GET for "/pastes/", a list of publicly available pastes.
func pastesGET(ec echo.Context) error {
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
