package storage

import "errors"

var (
	ErrValidation           = errors.New("validation error")
	ErrDSNEmpty             = errors.New("storage DSN is empty")
	ErrDSNMissingScheme     = errors.New("storage DSN missing scheme")
	ErrDSNUnsupportedScheme = errors.New("unsupported storage DSN scheme")
	ErrRedisInvalidDSN      = errors.New("invalid redis DSN")
)
