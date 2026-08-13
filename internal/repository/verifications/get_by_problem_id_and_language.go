package verifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *CodeVerificationRepository) GetByProblemIDAndLanguage(ctx context.Context, problemID string, language model.Language) (model.Verification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, problem_id, tests_code, created_at, submissions_count, lang
		FROM public.code_verification
		WHERE problem_id=$1 AND lang=$2;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, problemID, language.String())

	var record VerificationRecord
	if err := row.Scan(&record.ID, &record.ProblemID, &record.Code, &record.CreatedAt, &record.SubmissionsCount, &record.Language); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Verification{}, fmt.Errorf(
				"no verifications for problem_id=%s and language=%s found: %w",
				problemID,
				language.String(),
				errs.ErrNotFound,
			)
		}
		return model.Verification{}, err
	}

	return fromRecordToModel(record)
}
