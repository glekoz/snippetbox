package jwtAuth

import (
	"errors"
)

var (
	ErrUpdateJWTCookie     = errors.New("jwtAuth: need to update access token")
	ErrInvalidRefreshToken = errors.New("jwtAuth: invalid refresh token")
	ErrServerError         = errors.New("jwtAuth: server error")
)
