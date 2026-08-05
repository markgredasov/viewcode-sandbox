package sandbox

import "time"

type executionRequest struct {
	Language string
	Code     string
	Timeout  time.Duration
}

type executionResult struct {
	StdOut   string
	StdErr   string
	ExitCode int64
	Duration time.Duration
}
