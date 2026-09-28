package model

type TestsResult struct {
	TotalTestsNumber    int
	PassedTestsNumber   int
	FailedTestLog       string
	PassingDurationInMs float64
}

func NewTestsResult(totalTestsNumber, passedTestsNumber int, failedTestLog string, duration float64) TestsResult {
	return TestsResult{
		TotalTestsNumber:    totalTestsNumber,
		PassedTestsNumber:   passedTestsNumber,
		FailedTestLog:       failedTestLog,
		PassingDurationInMs: duration,
	}
}
