package model

import (
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type Complexity string

const (
	ComplexityEasy   = Complexity("easy")
	ComplexityMedium = Complexity("medium")
	ComplexityHard   = Complexity("hard")
)

func NewComplexity(complexity string) (Complexity, error) {
	complexityObjValue := Complexity(complexity)
	if !complexityObjValue.IsValid() {
		return "", fmt.Errorf("invalid complexity rate: %s: %w", complexity, errs.ErrInvalidArgument)
	}

	return complexityObjValue, nil
}

func (c Complexity) String() string {
	return string(c)
}

func (c Complexity) IsValid() bool {
	switch c {
	case ComplexityEasy, ComplexityMedium, ComplexityHard:
		return true
	}

	return false
}
