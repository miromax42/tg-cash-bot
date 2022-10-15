package currency

import "github.com/cockroachdb/errors"

var (
	ErrNotSuccess = errors.New("request should be succeed")
	ErrNotSupport = errors.New("currency is not supported")
)
