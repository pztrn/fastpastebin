package pastes

import (
	// stdlib
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	// local
	"go.dev.pztrn.name/fastpastebin/internal/captcha"
	"go.dev.pztrn.name/fastpastebin/internal/structs"
	"go.dev.pztrn.name/fastpastebin/internal/templater"

	// other
	"github.com/alecthomas/chroma/lexers"
	"github.com/labstack/echo"
)

// POST for "/paste/" which will create new paste and redirect to
// "/pastes/CREATED_PASTE_ID". This handler will do all the job for
// requests comes from browsers via web interface.
func pastePOSTWebInterface(ec echo.Context) error {
	// We should check if database connection available.
	dbConn := c.Database.GetDatabaseConnection()
	if c.Config.Database.Type != "flatfiles" && dbConn == nil {
		return ec.Redirect(http.StatusFound, "/database_not_available")
	}

	params, err := ec.FormParams()
	if err != nil {
		c.Logger.Error().Msg("Passed paste form is empty")

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

	if !strings.ContainsAny(params["paste-keep-for"][0], "Mmhd") && params["paste-keep-for"][0] != "forever" {
		c.Logger.Debug().Str("field value", params["paste-keep-for"][0]).Msg("'Keep paste for' field have invalid value")

		errtpl := templater.GetErrorTemplate(ec, "Invalid 'Paste should be available for' parameter passed. Please do not try to hack us ;).")

		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	// Verify captcha.
	if !captcha.Verify(params["paste-captcha-id"][0], params["paste-captcha-solution"][0]) {
		c.Logger.Debug().Str("captcha ID", params["paste-captcha-id"][0]).Str("captcha solution", params["paste-captcha-solution"][0]).Msg("Invalid captcha solution")

		errtpl := templater.GetErrorTemplate(ec, "Invalid captcha solution.")

		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	paste := &structs.Paste{
		Title:    params["paste-title"][0],
		Data:     params["paste-contents"][0],
		Language: params["paste-language"][0],
	}

	// Paste creation time in UTC.
	createdAt := time.Now().UTC()
	paste.CreatedAt = &createdAt

	// Parse "keep for" field.
	// Defaulting to "forever".
	keepFor := 0
	keepForUnit := 0

	if params["paste-keep-for"][0] != "forever" {
		keepForUnitRegex := regexp.MustCompile("[Mmhd]")

		keepForRaw := regexInts.FindAllString(params["paste-keep-for"][0], 1)[0]

		var err error

		keepFor, err = strconv.Atoi(keepForRaw)
		if err != nil {
			if params["paste-keep-for"][0] == "forever" {
				c.Logger.Debug().Msg("Keeping paste forever!")

				keepFor = 0
			} else {
				c.Logger.Debug().Err(err).Msg("Failed to parse 'Keep for' integer")

				errtpl := templater.GetErrorTemplate(ec, "Invalid 'Paste should be available for' parameter passed. Please do not try to hack us ;).")

				return ec.HTML(http.StatusBadRequest, errtpl)
			}
		}

		keepForUnitRaw := keepForUnitRegex.FindAllString(params["paste-keep-for"][0], 1)[0]
		keepForUnit = structs.PasteKeepsCorrelation[keepForUnitRaw]
	}

	paste.KeepFor = keepFor
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
		_ = paste.CreatePassword(pastePassword[0])
	}

	id, err2 := c.Database.SavePaste(paste)
	if err2 != nil {
		c.Logger.Error().Err(err2).Msg("Failed to save paste")

		errtpl := templater.GetErrorTemplate(ec, "Failed to save paste. Please, try again later.")

		return ec.HTML(http.StatusBadRequest, errtpl)
	}

	newPasteIDAsString := strconv.FormatInt(id, 10)
	c.Logger.Debug().Msg("Paste saved, URL: /paste/" + newPasteIDAsString)

	// Private pastes have it's timestamp in URL.
	if paste.Private {
		return ec.Redirect(http.StatusFound, "/paste/"+newPasteIDAsString+"/"+strconv.FormatInt(paste.CreatedAt.Unix(), 10))
	}

	return ec.Redirect(http.StatusFound, "/paste/"+newPasteIDAsString)
}
