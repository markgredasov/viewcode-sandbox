package errs

import "errors"

var (
	ErrInvalidLanguage = errors.New("this language is unavailable")
	ErrEmptyCode       = errors.New("code is empty")
	ErrTooLargeCode    = errors.New("code is too large")
)
