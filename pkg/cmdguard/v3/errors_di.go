package v3

import "errors"

// DI (dependency injection) sentinel errors.
var (
	// ErrInvalidScope indicates an invalid DI scope operation.
	ErrInvalidScope = errors.New("invalid scope")

	// ErrServiceNotFound indicates a service was not found in the DI container.
	ErrServiceNotFound = errors.New("service not found")

	// ErrServiceConstruction indicates a service provider failed during construction.
	ErrServiceConstruction = errors.New("service construction failed")

	// ErrServiceRegistration indicates a service failed to register in the DI container.
	ErrServiceRegistration = errors.New("service registration failed")
)
