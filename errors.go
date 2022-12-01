package watchdog

import "errors"

var (
	ErrEmptyResult = errors.New("empty result")
	Cancel         = errors.New("cancelled")
)
