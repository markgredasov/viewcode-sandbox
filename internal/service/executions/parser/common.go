package parser

import (
	"fmt"

	utils_string "github.com/markgredasov/viewcode-sandbox/internal/utils/string"
)

func parseDuration(line string) (float64, error) {
	durationInSec, err := utils_string.ParseRegexpToFloat64(line, `(\d+\.?\d*)s`)
	if err != nil {
		return 0.0, fmt.Errorf("duration: %w", err)
	}

	//nolint:mnd // not a magic number
	durationInMS := durationInSec * 1_000
	return durationInMS, nil
}
