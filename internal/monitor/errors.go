package monitor

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig       = errors.New("invalid config")
	ErrProbeInitialization = errors.New("probe initialization failed")
	ErrServiceNotFound     = errors.New("service not found")
)

// NewProbeInitializationError wraps a base error with additional context
func NewProbeInitializationError(serviceID, serviceName string, err error) error {
	return fmt.Errorf("%w: service %q (ID: %s): %w", ErrProbeInitialization, serviceName, serviceID, err)
}
