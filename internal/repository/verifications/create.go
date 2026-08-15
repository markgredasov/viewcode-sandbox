package verifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return model.Verification{}, fmt.Errorf(
					"verification code for problem_id=%s and language=%s already exists: %w",
					verification.ProblemID,
					verification.Language,
					errs.ErrAlreadyExists,
				)
			case "23503":
				return model.Verification{}, fmt.Errorf(
					"problem with id=%s not found: %w",
					verification.ProblemID,
					errs.ErrNotFound,
				)
			}
		}
		return model.Verification{}, err
	}

	return fromRecordToModel(record)
}
