package model

import (
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type Language string

const (
	LanguageGo     = Language("go")
	LanguagePython = Language("py")
	LanguageRust   = Language("rs")
)

func NewLanguage(lang string) (Language, error) {
	language := Language(lang)
	if !language.IsValid() {
		return "", fmt.Errorf("invalid language: %w", errs.ErrInvalidArgument)
	}

	return language, nil
}

func (l Language) String() string {
	return string(l)
}

func (l Language) IsValid() bool {
	switch l {
	case LanguageGo, LanguagePython, LanguageRust:
		return true
	}

	return false
}
