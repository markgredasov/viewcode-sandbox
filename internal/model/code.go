package model

import (
	"unicode/utf8"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

const (
	MaxCodeLength = 1000
)

type Code string

func NewCode(c string) (Code, error) {
	if utf8.RuneCountInString(c) == 0 {
		return Code(""), errs.ErrEmptyCode
	}

	if utf8.RuneCountInString(c) > MaxCodeLength {
		return Code(""), errs.ErrTooLargeCode
	}

	return Code(c), nil
}

func (c Code) String() string {
	return string(c)
}
