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
	// stdlib
	"os"
	"os/signal"
	"syscall"

	// local
	"gitlab.com/pztrn/fastpastebin/domains/indexpage"
	"gitlab.com/pztrn/fastpastebin/domains/pastes"
	"gitlab.com/pztrn/fastpastebin/internal/captcha"
	"gitlab.com/pztrn/fastpastebin/internal/context"
	"gitlab.com/pztrn/fastpastebin/internal/database"
	"gitlab.com/pztrn/fastpastebin/internal/templater"
)

func main() {
	c := context.New()
	c.Initialize()

	c.Logger.Info().Msg("Starting Fast Pastebin...")

	// Here goes initial initialization for packages that want CLI flags
	// to be added.

	// Parse flags.
	c.Flagger.Parse()

	// Continue loading.
	c.LoadConfiguration()
	c.InitializePost()
	database.New(c)
	c.Database.Initialize()
	templater.Initialize(c)

	captcha.New(c)

	indexpage.New(c)
	pastes.New(c)

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)
	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalHandler
		c.Shutdown()
		shutdownDone <- true
	}()

	<-shutdownDone
	os.Exit(0)
}
