package app

import "errors"

var (
	ErrAlreadyStarted    = errors.New("already started")
	ErrParameterRequired = errors.New("parameter required")
	ErrParseFailed       = errors.New("parse failed")
)
