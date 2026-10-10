package model

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

const (
	maxProblemNameLength = 70
	maxDescriptionLength = 10_000
)

type Problem struct {
	ProblemID      string
	Name           string
	Description    string
	CreatedAt      time.Time
	Languages      []Language
	ComplexityRate Complexity
	Categories     []Category
}

func NewProblem(
	problemID,
	name,
	description string,
	createdAt time.Time,
	languages []string,
	complexityRate string,
	categories []string,
) (Problem, error) {
	if utf8.RuneCountInString(name) <= 0 {
		return Problem{}, fmt.Errorf("empty problem name: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(name) >= maxProblemNameLength {
		return Problem{}, fmt.Errorf("problem name is too large: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) <= 0 {
		return Problem{}, fmt.Errorf("description is empty: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) >= maxDescriptionLength {
		return Problem{}, fmt.Errorf("description is too large: %w", errs.ErrInvalidArgument)
	}

	languageValueObjects, err := ConvertToValueObjsArray(languages, NewLanguage)
	if err != nil {
		return Problem{}, fmt.Errorf("convert languages: %w", err)
	}

	complexityRateValueObject, err := NewComplexity(complexityRate)
	if err != nil {
		return Problem{}, err
	}

	categoryValueObjects, err := ConvertToValueObjsArray(categories, NewCategory)
	if err != nil {
		return Problem{}, fmt.Errorf("convert categories: %w", err)
	}

	return Problem{
		ProblemID:      problemID,
		Name:           name,
		Description:    description,
		CreatedAt:      createdAt,
		Languages:      languageValueObjects,
		ComplexityRate: complexityRateValueObject,
		Categories:     categoryValueObjects,
	}, nil
}

func NewProblemUnitialized(name, description string) (Problem, error) {
	if utf8.RuneCountInString(name) <= 0 {
		return Problem{}, fmt.Errorf("empty problem name: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(name) >= maxProblemNameLength {
		return Problem{}, fmt.Errorf("problem name is too large: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) <= 0 {
		return Problem{}, fmt.Errorf("description is empty: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) >= maxDescriptionLength {
		return Problem{}, fmt.Errorf("description is too large: %w", errs.ErrInvalidArgument)
	}

	return Problem{
		ProblemID:   uuid.NewString(),
		Name:        name,
		Description: description,
	}, nil
}

func (p *Problem) WithComplexityRate(complexity Complexity) {
	p.ComplexityRate = complexity
}

func (p *Problem) WithCategories(categories []Category) {
	p.Categories = categories
}
