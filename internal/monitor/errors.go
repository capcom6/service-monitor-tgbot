package monitor

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig       = errors.New("invalid config")
	ErrProbeInitialization = errors.New("probe initialization failed")
)

// NewProbeInitializationError wraps a base error with additional context
func NewProbeInitializationError(serviceID, serviceName string, err error) error {
	return errors.Join(ErrProbeInitialization, fmt.Errorf("service %q (ID: %s): %w", serviceName, serviceID, err))
}
