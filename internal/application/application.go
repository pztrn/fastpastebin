package application

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	databaseinterface "go.dev.pztrn.name/fastpastebin/internal/database/interface"
)

// Application is a main application superstructure. It passes around all parts of application
// and serves as lifecycle management thing as well as kind-of-dependency-injection thing.
type Application struct {
	Config        *Config
	Database      databaseinterface.Interface
	Echo          *echo.Echo
	Log           zerolog.Logger
	services      map[string]Service
	servicesMutex sync.RWMutex
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

// New creates new application superstructure.
func New() *Application {
	//nolint:exhaustruct
	appl := &Application{}
	appl.initialize()

	return appl
}

// GetContext returns application-wide context.
func (a *Application) GetContext() context.Context {
	return a.ctx
}

// GetService returns interface{} with requested service or error if service wasn't registered.
func (a *Application) GetService(name string) (Service, error) {
	a.servicesMutex.RLock()
	srv, found := a.services[name]
	a.servicesMutex.RUnlock()

	if !found {
		return nil, fmt.Errorf("%s: %w", ErrApplicationError, ErrApplicationServiceNotRegistered)
	}

	return srv, nil
}

// Initializes internal state.
func (a *Application) initialize() {
	a.Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	a.Log.Info().Msg("Initializing Application...")

	a.ctx, a.cancelFunc = context.WithCancel(context.Background())

	a.services = make(map[string]Service)

	cfg, err := newConfig(a)
	if err != nil {
		a.Log.Fatal().Err(err).Msg("Failed to initialize configuration!")
	}

	a.Config = cfg

	a.initializeLogger()
	a.initializeHTTPServer()
}

// RegisterService registers service for later re-use everywhere it's needed.
func (a *Application) RegisterService(srv Service) error {
	a.servicesMutex.Lock()
	_, found := a.services[srv.GetName()]
	a.servicesMutex.Unlock()

	if found {
		return fmt.Errorf("%s: %w", ErrApplicationError, ErrApplicationServiceAlreadyRegistered)
	}

	if err := srv.Initialize(); err != nil {
		return fmt.Errorf("%s: %s: %w", ErrApplicationError, ErrApplicationServiceRegister, err)
	}

	a.services[srv.GetName()] = srv

	return nil
}

// Shutdown shutdowns application.
func (a *Application) Shutdown() error {
	a.cancelFunc()

	a.servicesMutex.RLock()
	defer a.servicesMutex.RUnlock()

	for _, service := range a.services {
		if err := service.Shutdown(); err != nil {
			return err
		}
	}

	return nil
}

// Start starts application.
func (a *Application) Start() error {
	a.initializeLoggerPost()
	a.startHTTPServer()

	a.servicesMutex.RLock()
	defer a.servicesMutex.RUnlock()

	for _, service := range a.services {
		if err := service.Start(); err != nil {
			return err
		}
	}

	return nil
}
