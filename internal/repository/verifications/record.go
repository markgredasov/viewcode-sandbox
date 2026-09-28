package verifications

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type VerificationRecord struct {
	ID           string
	ProblemID    string
	Code         string
	Language     string
	CreatedAt    time.Time
	TestsNumber  int
	CodeTemplate *string
}

func fromRecordToModel(r VerificationRecord) (model.Verification, error) {
	verificationModel, err := model.NewVerification(r.ID, r.ProblemID, r.Code, r.Language, r.CreatedAt, r.TestsNumber)
	if err != nil {
		return model.Verification{}, err
	}

	if err = verificationModel.WithCodeTemplate(r.CodeTemplate); err != nil {
		return model.Verification{}, err
	}

	return verificationModel, nil
}
