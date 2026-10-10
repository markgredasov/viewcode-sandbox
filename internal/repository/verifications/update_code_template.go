package verifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *CodeVerificationRepository) UpdateCodeTemplate(ctx context.Context, verificationID, codeTemplate string) (model.Verification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		UPDATE public.code_verification
		SET code_template = $1
		WHERE id = $2
		RETURNING id, problem_id, tests_code, created_at, total_number, lang, code_template;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, codeTemplate, verificationID)

	var record VerificationRecord
	err := row.Scan(&record.ID, &record.ProblemID, &record.Code, &record.CreatedAt, &record.TestsNumber, &record.Language, &record.CodeTemplate)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			msg := fmt.Sprintf("verification with id=%s not found", verificationID)
			return model.Verification{}, errs.NotFound(msg)
		}
		return model.Verification{}, errs.Internal("row scan", err)
	}

	return toModel(record)
}
