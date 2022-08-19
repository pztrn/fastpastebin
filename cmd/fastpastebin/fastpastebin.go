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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.dev.pztrn.name/fastpastebin/domains/dbnotavailable"
	"go.dev.pztrn.name/fastpastebin/domains/indexpage"
	"go.dev.pztrn.name/fastpastebin/domains/pastes"
	"go.dev.pztrn.name/fastpastebin/internal/application"
	"go.dev.pztrn.name/fastpastebin/internal/captcha"
	"go.dev.pztrn.name/fastpastebin/internal/database"
	"go.dev.pztrn.name/fastpastebin/internal/templater"
)

func main() {
	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)

	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)

	app := application.New()
	app.Log.Info().Msg("Starting Fast Pastebin...")

	database.New(app)
	app.Database.Initialize()
	templater.Initialize(app)

	captcha.New(app)

	dbnotavailable.New(app)
	indexpage.New(app)
	pastes.New(app)

	if err := app.Start(); err != nil {
		app.Log.Fatal().Err(err).Msg("Failed to start Fast Pastebin!")
	}

	go func() {
		<-signalHandler

		if err := app.Shutdown(); err != nil {
			app.Log.Error().Err(err).Msg("Fast Pastebin failed to shutdown!")
		}

		shutdownDone <- true
	}()

	<-shutdownDone
	os.Exit(0)
}
