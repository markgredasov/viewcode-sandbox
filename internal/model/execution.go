package model

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type Execution struct {
	ProblemID uuid.UUID

	Language Language
	Code     Code
	Tests    Verification

	ExecutionResult ExecutionResult
	TestsResult     TestsResult
}

func NewExecution(id, language, code string) (Execution, error) {
	uuidObject, err := uuid.Parse(id)
	if err != nil {
		return Execution{}, err
	}

	languageObject, err := NewLanguage(language)
	if err != nil {
		return Execution{}, err
	}

	codeObject, err := NewCode(code)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		ProblemID: uuidObject,
		Language:  languageObject,
		Code:      codeObject,
	}, nil
}

func (e *Execution) WithTestsResult(tr TestsResult) {
	e.TestsResult = tr
}

func (e *Execution) ConvertStdOutToTestResults() error {
	switch e.Language {
	case LanguageGo:
		e.WithTestsResult(getTestsResultFromStdOut(e.ExecutionResult.StdOut))
	case LanguagePython:

	case LanguageJava:
		fallthrough
	case LanguageCpp:
		fallthrough
	case LanguageJavaScript:
		fallthrough
	case LanguageTypeScript:
		fallthrough
	case LanguageCSharp:
		fallthrough
	case LanguageKotlin:
		fallthrough
	case LanguageSwift:
		fallthrough
	case LanguageRust:
		fallthrough
	default:
		return fmt.Errorf("language not provided yet: %w", errs.ErrNotImplemented)
	}

	return nil
}

func getTestsResultFromStdOut(stdOut string) TestsResult {
	scanner := bufio.NewScanner(strings.NewReader(stdOut))

	var testsCount int64 = 0
	var testsPassed int64 = 0
	lastLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			lastLine = line
		}

		var idx int
		if idx = strings.Index(line, "==="); idx != -1 {
			line = line[idx:]
		} else if idx = strings.Index(line, "---"); idx != -1 {
			line = line[idx:]
		}

		if strings.HasPrefix(line, "=== RUN") {
			testsCount++
		} else if strings.HasPrefix(line, "--- PASS") {
			testsPassed++
		}
	}

	duration := parseDuration(lastLine)

	return NewTestsResult(testsCount, testsPassed, "", duration)
}

func parseDuration(line string) float64 {
	re := regexp.MustCompile(`(\d+\.?\d*)s`)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 1 {
		durationInSec, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			//nolint:mnd // not a magic number
			durationInMS := durationInSec * 1_000
			return durationInMS
		}
	}

	return 0.0
}
