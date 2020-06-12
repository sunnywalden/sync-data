package errors

import "errors"

var (
	ErrRedisConfigNil = errors.New("redis config nil")
)
