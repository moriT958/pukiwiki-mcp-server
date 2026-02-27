package libpuki

import "errors"

var (
	ErrPageNotFound = errors.New("page not found")
	ErrAuthFailed   = errors.New("authentication failed")
)
