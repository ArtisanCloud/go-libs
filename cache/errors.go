package cache

import "errors"

var (
	ErrCacheMiss    = errors.New("cache: miss")
	ErrCASConflict  = errors.New("cache: compare-and-swap conflict")
	ErrNoStats      = errors.New("cache: no statistics available")
	ErrNotStored    = errors.New("cache: item not stored")
	ErrServerError  = errors.New("cache: server error")
	ErrInvalidValue = errors.New("cache: invalid value")
)
