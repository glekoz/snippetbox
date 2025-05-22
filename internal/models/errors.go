package models

import "errors"

var (
	ErrNoRecord         = errors.New("models: no matching record found")
	ErrWrongCredentials = errors.New("models: wrong credentials")
)
