package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrNotImplemented  = errors.New("not implemented yet")
	ErrInternal        = errors.New("internal")
	ErrTimeoutExceeded = errors.New("timeout exceeded")
)

func InvalidArgument(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrInvalidArgument)
}

func NotFound(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrNotFound)
}

func AlreadyExists(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrAlreadyExists)
}

func InvalidRequest(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrInvalidRequest)
}

func NotImplemented(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrNotImplemented)
}

func Internal(msg string, err error) error {
	return fmt.Errorf("%s: %w: %w", msg, err, ErrInternal)
}

func TimeoutExceeded(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrTimeoutExceeded)
}
