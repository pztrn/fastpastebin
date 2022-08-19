package application

import "errors"

var (
	// ErrApplicationError indicates that error belongs to Application.
	ErrApplicationError = errors.New("application")

	// ErrApplicationServiceRegister appears when trying to register (and initialize) a service
	// but something went wrong.
	ErrApplicationServiceRegister = errors.New("service registering and initialization")

	// ErrApplicationServiceAlreadyRegistered appears when trying to register service with already used name.
	ErrApplicationServiceAlreadyRegistered = errors.New("service already registered")

	// ErrApplicationServiceNotRegistered appears when trying to obtain a service that wasn't previously registered.
	ErrApplicationServiceNotRegistered = errors.New("service not registered")
)
