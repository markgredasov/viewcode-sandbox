package verification

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *VerificationRepository) GetByID(ctx context.Context, problemID string) (model.Verification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT tests_code
		FROM public.problems_tests
		WHERE problem_id = $1;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, problemID)

	var testsCode string
	if err := row.Scan(&testsCode); err != nil {
		return model.Verification{}, err
	}

	testsModel, err := model.NewVerification(testsCode)
	if err != nil {
		return model.Verification{}, err
	}

	return testsModel, nil
}
