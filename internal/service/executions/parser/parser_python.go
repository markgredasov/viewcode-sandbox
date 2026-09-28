package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	utils_string "github.com/markgredasov/viewcode-sandbox/internal/utils/string"
)

const (
	testsFailedRegexp = `(\d+)\s+failed`
	testsPassedRegexp = `(\d+)\s+passed`
)

type parserLanguagePython struct{}

func newParserLanguagePython() *parserLanguagePython {
	return &parserLanguagePython{}
}

func (p *parserLanguagePython) Parse(stdOut string) (model.TestsResult, error) {
	scanner := bufio.NewScanner(strings.NewReader(stdOut))

	lastLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			lastLine = line
		}
	}

	testsFailed, err := utils_string.ParseRegexpToInt(lastLine, testsFailedRegexp)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return model.TestsResult{}, fmt.Errorf("testsFailed: %w", err)
		}
	}

	testsPassed, err := utils_string.ParseRegexpToInt(lastLine, testsPassedRegexp)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return model.TestsResult{}, fmt.Errorf("testsPassed: %w", err)
		}
	}

	testsCount := testsPassed + testsFailed

	duration, err := parseDuration(lastLine)
	if err != nil {
		return model.TestsResult{}, fmt.Errorf("duration: %w", err)
	}

	return model.NewTestsResult(testsCount, testsPassed, "", duration), nil
}
