package parser

import (
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type registry struct {
	parsers map[model.Language]testsResultParser
}

func newRegistry() *registry {
	return &registry{
		parsers: make(map[model.Language]testsResultParser),
	}
}

func (r *registry) register(language model.Language, parser testsResultParser) error {
	if _, ok := r.parsers[language]; ok {
		return fmt.Errorf("parser: %w", errs.ErrAlreadyExists)
	}

	r.parsers[language] = parser
	return nil
}

func (r *registry) get(language model.Language) (testsResultParser, error) {
	parser, ok := r.parsers[language]
	if !ok {
		return nil, fmt.Errorf("parser: %w", errs.ErrNotFound)
	}

	return parser, nil
}

type testsResultParser interface {
	Parse(stdOut string) (model.TestsResult, error)
}
