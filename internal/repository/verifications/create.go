package verifications

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *CodeVerificationRepository) Create(ctx context.Context, verification model.Verification) (model.Verification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		INSERT INTO public.code_verification (id, problem_id, tests_code, lang)
		VALUES ($1, $2, $3, $4)
		RETURNING id, problem_id, tests_code, created_at, submissions_count, lang;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, verification.ID, verification.ProblemID, verification.Code, verification.Language)

	var record VerificationRecord
	err := row.Scan(&record.ID, &record.ProblemID, &record.Code, &record.CreatedAt, &record.SubmissionsCount, &record.Language)
	if err != nil {
		return model.Verification{}, err
	}

	return fromRecordToModel(record)
}
