package verifications

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type VerificationRecord struct {
	ID               string
	ProblemID        string
	Code             string
	Language         string
	CreatedAt        time.Time
	SubmissionsCount int64
	CodeTemplate     *string
}

func fromRecordToModel(r VerificationRecord) (model.Verification, error) {
	verificationModel, err := model.NewVerification(
		r.ID,
		r.ProblemID,
		r.Code,
		r.Language,
		r.CreatedAt,
		r.SubmissionsCount,
	)
	if err != nil {
		return model.Verification{}, err
	}

	if err = verificationModel.WithCodeTemplate(r.CodeTemplate); err != nil {
		return model.Verification{}, err
	}

	return verificationModel, nil
}
