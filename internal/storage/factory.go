package storage

import (
	"fmt"
	"net/url"
)

func NewFromDSN(dsn string) (Storage, error) {
	if dsn == "" {
		return nil, fmt.Errorf("%w", ErrDSNEmpty)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse storage DSN: %w", err)
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("%w: %s", ErrDSNMissingScheme, dsn)
	}

	switch u.Scheme {
	case "file":
		return newYamlStorage(u)
	case "redis":
		return newRedisStorage(u)
	default:
		return nil, fmt.Errorf("%w: %s", ErrDSNUnsupportedScheme, u.Scheme)
	}
}
