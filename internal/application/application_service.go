package application

// Service is a generic interface for all services of application.
type Service interface {
	// GetName returns service name for registering with application superstructure.
	GetName() string
	// Initialize initializes service.
	Initialize() error
	// Shutdown shuts service down if needed. Also should block is shutdown should be done in synchronous manner.
	Shutdown() error
	// Start starts service if needed. Should not block execution.
	Start() error
}
