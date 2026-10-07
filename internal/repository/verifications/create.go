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
		INSERT INTO public.code_verification (id, problem_id, tests_code, total_number, lang)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, problem_id, tests_code, created_at, total_number, lang;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, verification.ID, verification.ProblemID,
		verification.Code, verification.TestsNumber, verification.Language)

	var record VerificationRecord
	err := row.Scan(&record.ID, &record.ProblemID, &record.Code, &record.CreatedAt, &record.TestsNumber, &record.Language)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				msg := fmt.Sprintf("verification code for problem_id=%s and language=%s already exists",
					verification.ProblemID, verification.Language)
				return model.Verification{}, errs.AlreadyExists(msg)
			case "23503":
				msg := fmt.Sprintf("problem with id=%s not found", verification.ProblemID)
				return model.Verification{}, errs.NotFound(msg)
			}
		}
		return model.Verification{}, errs.Internal("row scan", err)
	}

	return toModel(record)
}
