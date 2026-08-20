package parser

import "github.com/markgredasov/viewcode-sandbox/internal/model"

type Adapter struct {
	registry *registry
}

func NewAdapter() (*Adapter, error) {
	registry := newRegistry()
	if err := registry.register(model.LanguagePython, newParserLanguagePython()); err != nil {
		return nil, err
	}

	if err := registry.register(model.LanguageGo, newParserLanguageGo()); err != nil {
		return nil, err
	}

	return &Adapter{
		registry: registry,
	}, nil
}

func (a *Adapter) Parse(language model.Language, stdOut string) (model.TestsResult, error) {
	testsResultParser, err := a.registry.get(language)
	if err != nil {
		return model.TestsResult{}, err
	}

	return testsResultParser.Parse(stdOut)
}
