package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	stringutils "github.com/markgredasov/viewcode-sandbox/internal/platform/stringutils"
)

const (
	testsFailedRegexp = `(\d+)\s+failed`
	testsOkRegexp     = `(\d+)\s+passed`
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

	testsFailed, err := stringutils.ParseRegexpToInt(lastLine, testsFailedRegexp)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return model.TestsResult{}, fmt.Errorf("testsFailed: %w", err)
		}
	}

	testsPassed, err := stringutils.ParseRegexpToInt(lastLine, testsOkRegexp)
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
