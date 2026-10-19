package verifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *CodeVerificationRepository) GetByProblemIDAndLanguage(ctx context.Context, problemID string, language model.Language) (model.Verification, error) {
	sqlQuery := `
		SELECT id, problem_id, tests_code, created_at, total_number, lang, code_template
		FROM public.code_verification
		WHERE problem_id=$1 AND lang=$2;
	`

	conn := dbpostgres.ExecutorFromContext(ctx, r.pool)

	rows, err := conn.Query(ctx, sqlQuery, problemID, language.String())
	if err != nil {
		return model.Verification{}, errs.Internal("rows query", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[VerificationRecord])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			msg := fmt.Sprintf("no verifications for problem_id=%s and language=%s found", problemID, language.String())
			return model.Verification{}, errs.NotFound(msg)
		}
		return model.Verification{}, errs.Internal("row scan", err)
	}

	return toModel(rec)
}
