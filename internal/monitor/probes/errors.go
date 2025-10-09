package probes

import (
	"errors"
	"fmt"
)

var (
	ErrProbeFailed         = errors.New("probe failed")
	ErrConfigError         = fmt.Errorf("%w: invalid config", ErrProbeFailed)
	ErrInfractructureError = fmt.Errorf("%w: infrastructure error", ErrProbeFailed)
	ErrResponseError       = fmt.Errorf("%w: response error", ErrProbeFailed)
)
