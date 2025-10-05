package telegram

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
	ErrTokenIsEmpty  = fmt.Errorf("%w: token is empty", ErrInvalidConfig)
)
