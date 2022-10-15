package repo

import "github.com/cockroachdb/errors"

var (
	ErrLimitExceed = errors.New("limit exceed")
)
