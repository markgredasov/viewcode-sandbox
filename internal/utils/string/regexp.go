package utils_string

import (
	"fmt"
	"regexp"
	"strconv"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

func ParseRegexpToInt(line string, regExpression string) (int, error) {
	re := regexp.MustCompile(regExpression)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 1 {
		result, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("convert: %w", err)
		}
		return result, nil
	}

	return 0, fmt.Errorf("no matches: %w", errs.ErrNotFound)
}

func ParseRegexpToFloat64(line string, regExpression string) (float64, error) {
	re := regexp.MustCompile(regExpression)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 1 {
		result, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0.0, fmt.Errorf("convert: %w", err)
		}
		return result, nil
	}

	return 0.0, fmt.Errorf("no matches: %w", errs.ErrNotFound)
}
