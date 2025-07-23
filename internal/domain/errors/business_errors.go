package businesserrors

import "errors"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrNotFound           = errors.New("not found")
	ErrSomethingWentWrong = errors.New("something went wrong")
)
