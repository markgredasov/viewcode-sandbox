package verifications

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type VerificationRecord struct {
	ID           string    `db:"id"`
	ProblemID    string    `db:"problem_id"`
	Code         string    `db:"tests_code"`
	Language     string    `db:"lang"`
	CreatedAt    time.Time `db:"created_at"`
	TestsNumber  int       `db:"total_number"`
	CodeTemplate *string   `db:"code_template"`
}

func toModel(r VerificationRecord) (model.Verification, error) {
	verificationModel, err := model.NewVerification(r.ID, r.ProblemID, r.Code, r.Language, r.CreatedAt, r.TestsNumber)
	if err != nil {
		return model.Verification{}, err
	}

	if err = verificationModel.WithCodeTemplate(r.CodeTemplate); err != nil {
		return model.Verification{}, err
	}

	return verificationModel, nil
}
