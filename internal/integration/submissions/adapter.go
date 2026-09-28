package submissionsingegration

import (
	"context"
	"fmt"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
)

type Adapter struct {
	client *client
}

func NewAdapter(cfg Config) (*Adapter, error) {
	cl, err := newClient(cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("create adapter: %w", err)
	}

	return &Adapter{
		client: cl,
	}, nil
}

func (a *Adapter) CreateSubmission(ctx context.Context, userID string, problemID string,
	testsTotal int, language model.Language) (submission.Submission, error) {
	return submission.Submission{}, nil
}
func (a *Adapter) SetRunning(ctx context.Context, id string) (submission.Submission, error) {
	return submission.Submission{}, nil
}
func (a *Adapter) SetResult(ctx context.Context, id string, status submission.Status,
	testsPassed int, executionTimeInMS float64) (submission.Submission, error) {
	return submission.Submission{}, nil
}
