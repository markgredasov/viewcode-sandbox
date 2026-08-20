package parser

import (
	"bufio"
	"errors"
	"strings"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

const (
	testRunSubstring    = "=== RUN"
	testPassedSubstring = "--- PASS"
)

type parserLanguageGo struct{}

func newParserLanguageGo() *parserLanguageGo {
	return &parserLanguageGo{}
}

func (p *parserLanguageGo) Parse(stdOut string) (model.TestsResult, error) {
	scanner := bufio.NewScanner(strings.NewReader(stdOut))

	var testsCount int64 = 0
	var testsPassed int64 = 0
	lastLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			lastLine = line
		}

		if strings.Contains(line, testRunSubstring) {
			testsCount++
		} else if strings.Contains(line, testPassedSubstring) {
			testsPassed++
		}
	}

	duration, err := parseDuration(lastLine)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return model.TestsResult{}, err
		}
	}

	return model.NewTestsResult(testsCount, testsPassed, "", duration), nil
}
