package singleton

import "errors"

var (
	ErrAlreadyRunning = errors.New("already running")
)
