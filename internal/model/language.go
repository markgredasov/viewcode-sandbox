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
		return "", fmt.Errorf("invalid language: %s: %w", lang, errs.ErrInvalidArgument)
	}

	return language, nil
}

func ConvertTo(languages []string) ([]Language, error) {
	languageValueObjects := make([]Language, len(languages))
	for i := range languages {
		languageValueObject, err := NewLanguage(languages[i])
		if err != nil {
			return []Language{}, err
		}

		languageValueObjects[i] = languageValueObject
	}

	return languageValueObjects, nil
}

func ConvertFrom(languages []Language) []string {
	languagesStr := make([]string, len(languages))
	for i := range languages {
		languagesStr[i] = languages[i].String()
	}

	return languagesStr
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
