package verifications

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *CodeVerificationRepository) GetByID(ctx context.Context, problemID string) (model.Verification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT tests_code
		FROM public.code_verification
		WHERE problem_id = $1;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, problemID)

	var testsCode string
	if err := row.Scan(&testsCode); err != nil {
		return model.Verification{}, err
	}

	testsModel, err := model.NewVerificationByCodeValue(testsCode)
	if err != nil {
		return model.Verification{}, err
	}

	return testsModel, nil
}
