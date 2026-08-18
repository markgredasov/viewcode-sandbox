package service_executions

import (
	"bufio"
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ExecutionsService) executeLanguageGo(ctx context.Context, exec model.Execution, tests model.Verification) (model.Execution, error) {
	execution, err := s.testsRunningPort.RunTests(ctx, exec, tests)
	if err != nil {
		return model.Execution{}, err
	}

	execution.WithTestsResult(getTestsResultFromStdOut(execution.ExecutionResult.StdOut))

	return execution, nil
}

func getTestsResultFromStdOut(stdOut string) model.TestsResult {
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

	return model.NewTestsResult(testsCount, testsPassed, "", duration)
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
