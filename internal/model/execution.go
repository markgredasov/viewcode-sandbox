package model

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	utils_string "github.com/markgredasov/viewcode-sandbox/internal/utils/string"
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
		testsResult, err := getTestsResultFromStdOutLanguageGo(e.ExecutionResult.StdOut)
		if err != nil {
			return err
		}
		e.WithTestsResult(testsResult)
	case LanguagePython:
		testsResult, err := getTestsResultFromStdOutLanguagePython(e.ExecutionResult.StdOut)
		if err != nil {
			return err
		}
		e.WithTestsResult(testsResult)
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

func getTestsResultFromStdOutLanguagePython(stdOut string) (TestsResult, error) {
	scanner := bufio.NewScanner(strings.NewReader(stdOut))

	lastLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			lastLine = line
		}
	}

	testsFailed, err := utils_string.ParseRegexpToInt64(lastLine, `(\d+)\s+failed`)
	if err != nil {
		//nolint:govet // will change this later
		if !errors.As(err, &errs.ErrNotFound) {
			return TestsResult{}, fmt.Errorf("testsFailed: %w", err)
		}
	}

	testsPassed, err := utils_string.ParseRegexpToInt64(lastLine, `(\d+)\s+passed`)
	if err != nil {
		return TestsResult{}, fmt.Errorf("testsPassed: %w", err)
	}

	testsCount := testsPassed + testsFailed

	duration, err := parseDuration(lastLine)
	if err != nil {
		return TestsResult{}, fmt.Errorf("duration: %w", err)
	}

	return NewTestsResult(testsCount, testsPassed, "", duration), nil
}

func getTestsResultFromStdOutLanguageGo(stdOut string) (TestsResult, error) {
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

	duration, err := parseDuration(lastLine)
	if err != nil {
		return TestsResult{}, err
	}

	return NewTestsResult(testsCount, testsPassed, "", duration), nil
}

func parseDuration(line string) (float64, error) {
	durationInSec, err := utils_string.ParseRegexpToFloat64(line, `(\d+\.?\d*)s`)
	if err != nil {
		return 0.0, err
	}

	//nolint:mnd // not a magic number
	durationInMS := durationInSec * 1_000
	return durationInMS, nil
}
