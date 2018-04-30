package main

import (
	// stdlib
	"os"
	"os/signal"
	"syscall"

	// local
	"github.com/pztrn/fastpastebin/api"
	"github.com/pztrn/fastpastebin/context"
	"github.com/pztrn/fastpastebin/database"
	"github.com/pztrn/fastpastebin/database/migrations"
	"github.com/pztrn/fastpastebin/pastes"
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
	database.New(c)
	c.Database.Initialize()
	migrations.New(c)
	migrations.Migrate()
	api.New(c)
	api.InitializeAPI()

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
