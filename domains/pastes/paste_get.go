package pastes

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/alecthomas/chroma"
	htmlfmt "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/labstack/echo"
	"go.dev.pztrn.name/fastpastebin/internal/database/dialects/flatfiles"
	"go.dev.pztrn.name/fastpastebin/internal/structs"
	"go.dev.pztrn.name/fastpastebin/internal/templater"
)

const (
	pasteCookieInvalid    = "PASTE_COOKIE_INVALID"
	pasteExpired          = "PASTE_EXPIRED"
	pasteNotFound         = "PASTE_NOT_FOUND"
	pasteTimestampInvalid = "PASTE_TIMESTAMP_INVALID"
)

// Actual getting paste data and returns it's content without formatting.
// This function will return paste's structure and optional error string
// that defined in constants above.
// Actually required only paste ID, all other parameters are optional
// for some cases, e.g. public paste won't check for timestamp and cookie
// value (they both will be ignored), but private will.
func pasteGetData(pasteID int, timestamp int64, cookieValue string) (*structs.Paste, string) {
	// Get paste.
	paste, err1 := ctx.Database.GetPaste(pasteID)
	if err1 != nil {
		ctx.Logger.Error().Err(err1).Int("paste ID", pasteID).Msg("Failed to get paste")

		return nil, pasteNotFound
	}

	// Check if paste is expired.
	if paste.IsExpired() {
		ctx.Logger.Error().Int("paste ID", pasteID).Msg("Paste is expired")

		return nil, pasteExpired
	}

	// Check if we have a private paste and it's parameters are correct.
	if paste.Private {
		pasteTS := paste.CreatedAt.Unix()
		if timestamp != pasteTS {
			ctx.Logger.Error().Int("paste ID", pasteID).Int64("paste timestamp", pasteTS).Int64("provided timestamp", timestamp).Msg("Incorrect timestamp provided for private paste")

			return nil, pasteTimestampInvalid
		}
	}

	// If we have a private paste requested and password for that paste
	// was defined - check additional things that required to view this
	// paste.
	if paste.Private && paste.Password != "" {
		// Generate cookie value to check.
		pasteCookieValue := paste.GenerateCryptedCookieValue()

		if cookieValue != pasteCookieValue {
			return nil, pasteCookieInvalid
		}
	}

	return paste, ""
}

// GET for "/paste/PASTE_ID" and "/paste/PASTE_ID/TIMESTAMP" (private pastes).
// Web interface version.
func pasteGETWebInterface(ectx echo.Context) error {
	pasteIDRaw := ectx.Param("id")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	pasteIDStr := strconv.Itoa(pasteID)
	ctx.Logger.Debug().Int("paste ID", pasteID).Msg("Trying to get paste data")

	// Check if we have timestamp passed.
	// If passed timestamp is invalid (isn't a real UNIX timestamp) we
	// will show 404 Not Found error and spam about that in logs.
	var timestamp int64

	tsProvidedStr := ectx.Param("timestamp")
	if tsProvidedStr != "" {
		tsProvided, err := strconv.ParseInt(tsProvidedStr, 10, 64)
		if err != nil {
			ctx.Logger.Error().Err(err).Int("paste ID", pasteID).Int64("provided timestamp", tsProvided).Msg("Invalid timestamp provided for getting private paste")

			errtpl := templater.GetErrorTemplate(ectx, "Paste #"+pasteIDStr+" not found")

			// nolint:wrapcheck
			return ectx.HTML(http.StatusBadRequest, errtpl)
		}

		timestamp = tsProvided
	}

	// Check if we have "PASTE-PASTEID" cookie defined. It is required
	// for private pastes.
	var cookieValue string

	cookie, err1 := ectx.Cookie("PASTE-" + pasteIDStr)
	if err1 == nil {
		cookieValue = cookie.Value
	}

	paste, err := pasteGetData(pasteID, timestamp, cookieValue)

	// For these cases we should return 404 Not Found page.
	if err == pasteExpired || err == pasteNotFound || err == pasteTimestampInvalid {
		errtpl := templater.GetErrorTemplate(ectx, "Paste #"+pasteIDRaw+" not found")

		// nolint:wrapcheck
		return ectx.HTML(http.StatusNotFound, errtpl)
	}

	// If passed cookie value was invalid - go to paste authorization
	// page.
	if err == pasteCookieInvalid {
		ctx.Logger.Info().Int("paste ID", pasteID).Msg("Invalid cookie, redirecting to auth page")

		// nolint:wrapcheck
		return ectx.Redirect(http.StatusMovedPermanently, "/paste/"+pasteIDStr+"/"+ectx.Param("timestamp")+"/verify")
	}

	// Format paste data map.
	pasteData := make(map[string]string)
	pasteData["pasteTitle"] = paste.Title
	pasteData["pasteID"] = strconv.Itoa(paste.ID)
	pasteData["pasteDate"] = paste.CreatedAt.Format("2006-01-02 @ 15:04:05") + " UTC"
	pasteData["pasteLanguage"] = paste.Language

	pasteExpirationString := "Never"
	if paste.KeepFor != 0 && paste.KeepForUnitType != 0 {
		pasteExpirationString = paste.GetExpirationTime().Format("2006-01-02 @ 15:04:05") + " UTC"
	}

	pasteData["pasteExpiration"] = pasteExpirationString

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
		ctx.Logger.Error().Err(err3).Msg("Failed to tokenize paste data")
	}
	// Get style for HTML output.
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	// Get HTML formatter.
	formatter := chroma.Formatter(htmlfmt.New(htmlfmt.WithLineNumbers(true), htmlfmt.LineNumbersInTable(true), htmlfmt.LinkableLineNumbers(true, "L")))

	// Create buffer and format into it.
	buf := new(bytes.Buffer)

	err4 := formatter.Format(buf, style, lexered)
	if err4 != nil {
		ctx.Logger.Error().Err(err4).Msg("Failed to format paste data")
	}

	pasteData["pastedata"] = buf.String()

	// Get template and format it.
	pasteHTML := templater.GetTemplate(ectx, "paste.html", pasteData)

	// nolint:wrapcheck
	return ectx.HTML(http.StatusOK, pasteHTML)
}

// GET for "/paste/PASTE_ID/TIMESTAMP/verify" - a password verify page.
func pastePasswordedVerifyGet(ectx echo.Context) error {
	pasteIDRaw := ectx.Param("id")
	timestampRaw := ectx.Param("timestamp")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])

	// Get paste.
	paste, err1 := ctx.Database.GetPaste(pasteID)
	if err1 != nil {
		ctx.Logger.Error().Err(err1).Int("paste ID", pasteID).Msg("Failed to get paste data")

		errtpl := templater.GetErrorTemplate(ectx, "Paste #"+pasteIDRaw+" not found")

		// nolint:wrapcheck
		return ectx.HTML(http.StatusBadRequest, errtpl)
	}

	// Check for auth cookie. If present - redirect to paste.
	cookie, err := ectx.Cookie("PASTE-" + strconv.Itoa(pasteID))
	if err == nil {
		// No cookie, redirect to auth page.
		ctx.Logger.Debug().Msg("Paste cookie found, checking it...")

		// Generate cookie value to check.
		cookieValue := paste.GenerateCryptedCookieValue()

		if cookieValue == cookie.Value {
			ctx.Logger.Info().Msg("Valid cookie, redirecting to paste page...")

			// nolint:wrapcheck
			return ectx.Redirect(http.StatusMovedPermanently, "/paste/"+pasteIDRaw+"/"+ectx.Param("timestamp"))
		}

		ctx.Logger.Debug().Msg("Invalid cookie, showing auth page")
	}

	// HTML data.
	htmlData := make(map[string]string)
	htmlData["pasteID"] = strconv.Itoa(pasteID)
	htmlData["pasteTimestamp"] = timestampRaw

	verifyHTML := templater.GetTemplate(ectx, "passworded_paste_verify.html", htmlData)

	// nolint:wrapcheck
	return ectx.HTML(http.StatusOK, verifyHTML)
}

// POST for "/paste/PASTE_ID/TIMESTAMP/verify" - a password verify page.
func pastePasswordedVerifyPost(ectx echo.Context) error {
	// We should check if database connection available.
	dbConn := ctx.Database.GetDatabaseConnection()

	if ctx.Config.Database.Type != flatfiles.FlatFileDialect && dbConn == nil {
		// nolint:wrapcheck
		return ectx.Redirect(http.StatusFound, "/database_not_available")
	}

	pasteIDRaw := ectx.Param("id")
	timestampRaw := ectx.Param("timestamp")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	ctx.Logger.Debug().Int("paste ID", pasteID).Msg("Requesting paste")

	// Get paste.
	paste, err1 := ctx.Database.GetPaste(pasteID)
	if err1 != nil {
		ctx.Logger.Error().Err(err1).Int("paste ID", pasteID).Msg("Failed to get paste")
		errtpl := templater.GetErrorTemplate(ectx, "Paste #"+strconv.Itoa(pasteID)+" not found")

		// nolint:wrapcheck
		return ectx.HTML(http.StatusBadRequest, errtpl)
	}

	params, err2 := ectx.FormParams()
	if err2 != nil {
		ctx.Logger.Debug().Msg("No form parameters passed")

		errtpl := templater.GetErrorTemplate(ectx, "Paste #"+strconv.Itoa(pasteID)+" not found")

		// nolint:wrapcheck
		return ectx.HTML(http.StatusBadRequest, errtpl)
	}

	if paste.VerifyPassword(params["paste-password"][0]) {
		// Set cookie that this paste's password is verified and paste
		// can be viewed.
		cookie := new(http.Cookie)
		cookie.Name = "PASTE-" + strconv.Itoa(pasteID)
		cookie.Value = paste.GenerateCryptedCookieValue()
		cookie.Expires = time.Now().Add(24 * time.Hour)
		ectx.SetCookie(cookie)

		// nolint:wrapcheck
		return ectx.Redirect(http.StatusFound, "/paste/"+strconv.Itoa(pasteID)+"/"+timestampRaw)
	}

	errtpl := templater.GetErrorTemplate(ectx, "Invalid password. Please, try again.")

	// nolint:wrapcheck
	return ectx.HTML(http.StatusBadRequest, errtpl)
}

// GET for "/pastes/:id/raw", raw paste output.
// Web interface version.
func pasteRawGETWebInterface(ectx echo.Context) error {
	// We should check if database connection available.
	dbConn := ctx.Database.GetDatabaseConnection()
	if ctx.Config.Database.Type != flatfiles.FlatFileDialect && dbConn == nil {
		// nolint:wrapcheck
		return ectx.Redirect(http.StatusFound, "/database_not_available/raw")
	}

	pasteIDRaw := ectx.Param("id")
	// We already get numbers from string, so we will not check strconv.Atoi()
	// error.
	pasteID, _ := strconv.Atoi(regexInts.FindAllString(pasteIDRaw, 1)[0])
	ctx.Logger.Debug().Int("paste ID", pasteID).Msg("Requesting paste data")

	// Get paste.
	paste, err1 := ctx.Database.GetPaste(pasteID)
	if err1 != nil {
		ctx.Logger.Error().Err(err1).Int("paste ID", pasteID).Msg("Failed to get paste from database")

		// nolint:wrapcheck
		return ectx.HTML(http.StatusBadRequest, "Paste #"+pasteIDRaw+" does not exist.")
	}

	if paste.IsExpired() {
		ctx.Logger.Error().Int("paste ID", pasteID).Msg("Paste is expired")

		// nolint:wrapcheck
		return ectx.HTML(http.StatusBadRequest, "Paste #"+pasteIDRaw+" does not exist.")
	}

	// Check if we have a private paste and it's parameters are correct.
	if paste.Private {
		tsProvidedStr := ectx.Param("timestamp")

		tsProvided, err2 := strconv.ParseInt(tsProvidedStr, 10, 64)
		if err2 != nil {
			ctx.Logger.Error().Err(err2).Int("paste ID", pasteID).Str("provided timestamp", tsProvidedStr).Msg("Invalid timestamp provided for getting private paste")

			// nolint:wrapcheck
			return ectx.String(http.StatusBadRequest, "Paste #"+pasteIDRaw+" not found")
		}

		pasteTS := paste.CreatedAt.Unix()
		if tsProvided != pasteTS {
			ctx.Logger.Error().Int("paste ID", pasteID).Int64("provided timestamp", tsProvided).Int64("paste timestamp", pasteTS).Msg("Incorrect timestamp provided for private paste")

			// nolint:wrapcheck
			return ectx.String(http.StatusBadRequest, "Paste #"+pasteIDRaw+" not found")
		}
	}

	// nolint
	// ToDo: figure out how to handle passworded pastes here.
	// Return error for now.
	if paste.Password != "" {
		ctx.Logger.Error().Int("paste ID", pasteID).Msg("Cannot render paste as raw: passworded paste. Patches welcome!")
		return ectx.String(http.StatusBadRequest, "Paste #"+pasteIDRaw+" not found")
	}

	// nolint:wrapcheck
	return ectx.String(http.StatusOK, paste.Data)
}
