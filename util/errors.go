package util

import "errors"

var (
	ErrBadFormat   = errors.New("bad format")
	ErrUnsupported = errors.New("not supported")
	ErrBadResponse = errors.New("bad response")
)
