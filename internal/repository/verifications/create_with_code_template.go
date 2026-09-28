package verifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *CodeVerificationRepository) CreateWithCodeTemplate(ctx context.Context, verification model.Verification) (model.Verification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		INSERT INTO public.code_verification (id, problem_id, tests_code, total_number, lang, code_template)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, problem_id, tests_code, created_at, total_number, lang, code_template;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, verification.ID, verification.ProblemID, verification.Code,
		verification.TestsNumber, verification.Language, verification.CodeTemplate)

	var record VerificationRecord
	err := row.Scan(&record.ID, &record.ProblemID, &record.Code, &record.CreatedAt, &record.TestsNumber, &record.Language, &record.CodeTemplate)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return model.Verification{}, fmt.Errorf(
				"problem with id=%s not found: %w",
				verification.ProblemID,
				errs.ErrNotFound,
			)
		}
		return model.Verification{}, err
	}

	return fromRecordToModel(record)
}
