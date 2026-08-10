package model

import (
	"fmt"
	"unicode/utf8"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

const (
	MaxCodeLength = 10000
)

type Code string

func NewCode(c string) (Code, error) {
	if utf8.RuneCountInString(c) == 0 {
		return Code(""), fmt.Errorf("code is empty: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(c) > MaxCodeLength {
		return Code(""), fmt.Errorf("too large code: %w", errs.ErrInvalidArgument)
	}

	return Code(c), nil
}

func (c Code) String() string {
	return string(c)
}
