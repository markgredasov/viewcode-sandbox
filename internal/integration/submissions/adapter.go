package submissionsingegration

import (
	"context"
	"fmt"
	"math"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
	submissions_v1 "github.com/markgredasov/viewcode-submissions/pkg/api/v1"
)

var (
	languageModelToGRPCModel = map[model.Language]submissions_v1.Language{
		model.LanguageGo:         submissions_v1.Language_LANGUAGE_GO,
		model.LanguagePython:     submissions_v1.Language_LANGUAGE_PYTHON,
		model.LanguageJava:       submissions_v1.Language_LANGUAGE_JAVA,
		model.LanguageCpp:        submissions_v1.Language_LANGUAGE_CPP,
		model.LanguageJavaScript: submissions_v1.Language_LANGUAGE_JAVA_SCRIPT,
		model.LanguageTypeScript: submissions_v1.Language_LANGUAGE_TYPE_SCRIPT,
		model.LanguageCSharp:     submissions_v1.Language_LANGUAGE_CSHARP,
		model.LanguageKotlin:     submissions_v1.Language_LANGUAGE_KOTLIN,
		model.LanguageSwift:      submissions_v1.Language_LANGUAGE_SWIFT,
		model.LanguageRust:       submissions_v1.Language_LANGUAGE_RUST,
	}

	//nolint:exhaustive // unspecified not needed in map
	languageGRPCModelToModel = map[submissions_v1.Language]model.Language{
		submissions_v1.Language_LANGUAGE_GO:          model.LanguageGo,
		submissions_v1.Language_LANGUAGE_PYTHON:      model.LanguagePython,
		submissions_v1.Language_LANGUAGE_JAVA:        model.LanguageJava,
		submissions_v1.Language_LANGUAGE_CPP:         model.LanguageCpp,
		submissions_v1.Language_LANGUAGE_JAVA_SCRIPT: model.LanguageJavaScript,
		submissions_v1.Language_LANGUAGE_TYPE_SCRIPT: model.LanguageTypeScript,
		submissions_v1.Language_LANGUAGE_CSHARP:      model.LanguageCSharp,
		submissions_v1.Language_LANGUAGE_KOTLIN:      model.LanguageKotlin,
		submissions_v1.Language_LANGUAGE_SWIFT:       model.LanguageSwift,
		submissions_v1.Language_LANGUAGE_RUST:        model.LanguageRust,
	}

	statusModelToGRPCModel = map[submission.Status]submissions_v1.SubmissionStatus{
		submission.StatusPending:     submissions_v1.SubmissionStatus_STATUS_PENDING,
		submission.StatusRunning:     submissions_v1.SubmissionStatus_STATUS_RUNNING,
		submission.StatusAccepted:    submissions_v1.SubmissionStatus_STATUS_ACCEPTED,
		submission.StatusWrongAnswer: submissions_v1.SubmissionStatus_STATUS_WRONG_ANSWER,
		submission.StatusTimeout:     submissions_v1.SubmissionStatus_STATUS_TIMEOUT,
		submission.StatusInternal:    submissions_v1.SubmissionStatus_STATUS_INTERNAL,
	}

	//nolint:exhaustive // unspecified not needed in map
	statusGRPCModelToModel = map[submissions_v1.SubmissionStatus]submission.Status{
		submissions_v1.SubmissionStatus_STATUS_PENDING:      submission.StatusPending,
		submissions_v1.SubmissionStatus_STATUS_RUNNING:      submission.StatusRunning,
		submissions_v1.SubmissionStatus_STATUS_ACCEPTED:     submission.StatusAccepted,
		submissions_v1.SubmissionStatus_STATUS_WRONG_ANSWER: submission.StatusWrongAnswer,
		submissions_v1.SubmissionStatus_STATUS_TIMEOUT:      submission.StatusTimeout,
		submissions_v1.SubmissionStatus_STATUS_INTERNAL:     submission.StatusInternal,
	}
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

func (a *Adapter) CreateSubmission(ctx context.Context, userID, problemID string,
	testsTotal int, language model.Language,
) (submission.Submission, error) {
	lang, ok := languageModelToGRPCModel[language]
	if !ok {
		return submission.Submission{}, errs.InvalidRequest(fmt.Sprintf("convert language to grpc model: %s", language))
	}

	if testsTotal < math.MinInt32 || testsTotal > math.MaxInt32 {
		return submission.Submission{}, errs.InvalidRequest(fmt.Sprintf("cannot convert testsTotal=%d to int32", testsTotal))
	}

	request := &submissions_v1.CreateRequest{
		UserId:     userID,
		ProblemId:  problemID,
		TestsTotal: int32(testsTotal),
		Language:   lang,
	}

	resp, err := a.client.CreateSubmission(ctx, request)
	if err != nil {
		return submission.Submission{}, err
	}

	return toSubmission(resp.GetSubmission())
}

func (a *Adapter) SetRunning(ctx context.Context, id string) (submission.Submission, error) {
	request := &submissions_v1.SetRunningRequest{
		Id: id,
	}

	resp, err := a.client.SetRunning(ctx, request)
	if err != nil {
		return submission.Submission{}, err
	}

	return toSubmission(resp.GetSubmission())
}

func (a *Adapter) SetResult(ctx context.Context, id string, status submission.Status,
	testsPassed int, executionTimeInMS float64,
) (submission.Submission, error) {
	stat, ok := statusModelToGRPCModel[status]
	if !ok {
		return submission.Submission{}, errs.InvalidRequest(fmt.Sprintf("convert status to grpc model: %s", status.String()))
	}

	if testsPassed < math.MinInt32 || testsPassed > math.MaxInt32 {
		return submission.Submission{}, errs.InvalidRequest(fmt.Sprintf("cannot convert testsPassed=%d to int32", testsPassed))
	}

	request := &submissions_v1.SetResultRequest{
		Id:              id,
		TestsPassed:     int32(testsPassed),
		Status:          stat,
		ExecutionTimeMs: int32(executionTimeInMS),
	}

	resp, err := a.client.SetResult(ctx, request)
	if err != nil {
		return submission.Submission{}, err
	}

	return toSubmission(resp.GetSubmission())
}

func toSubmission(from *submissions_v1.Submission) (submission.Submission, error) {
	lang, ok := languageGRPCModelToModel[from.GetLanguage()]
	if !ok {
		err := errs.ErrInternal
		return submission.Submission{}, errs.Internal(fmt.Sprintf("convert language from grpc model: %s", from.GetLanguage().String()), err)
	}

	status, ok := statusGRPCModelToModel[from.GetStatus()]
	if !ok {
		err := errs.ErrInternal
		return submission.Submission{}, errs.Internal(fmt.Sprintf("convert status from grpc model: %s", from.GetStatus().String()), err)
	}

	result, err := submission.New(from.GetUserId(), from.GetProblemId(), int(from.GetTestsTotal()),
		int(from.GetTestsPassed()), lang, int(from.GetExecutionTimeMs()))
	if err != nil {
		return submission.Submission{}, err
	}

	id, err := model.ToUUID(from.GetId())
	if err != nil {
		return submission.Submission{}, err
	}

	return *result.WithID(id).WithCreatedAt(from.GetCreatedAt().AsTime()).WithUpdatedAt(from.GetUpdatedAt().AsTime()).WithStatus(status), nil
}
