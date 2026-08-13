package errs

import "errors"

var (
	ErrNoTestsProvided = errors.New("no tests found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
)
