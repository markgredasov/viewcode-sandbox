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
}

func fromRecordToModel(r VerificationRecord) (model.Verification, error) {
	return model.NewVerification(
		r.ID,
		r.ProblemID,
		r.Code,
		r.Language,
		r.CreatedAt,
		r.SubmissionsCount,
	)
}
