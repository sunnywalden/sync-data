package errs

import "errors"

var (
	ErrRedisConfigNil = errors.New("redis config nil")
	ErrUserNotExists = errors.New("no user matched exists")
	ErrQueryParamsNil = errors.New("query params are all nil")
	ErrNotAuthorized = errors.New("not authorized")
	ErrTokenNotExists = errors.New("token could not be nil")
)
