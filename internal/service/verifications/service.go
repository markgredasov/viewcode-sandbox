package service_verifications

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type VerificationsService struct {
	verificationsRepository VerificationsRepository
}

type VerificationsRepository interface {
	Create(ctx context.Context, verification model.Verification) (model.Verification, error)
}

func NewVerificationsService(repo VerificationsRepository) *VerificationsService {
	return &VerificationsService{
		verificationsRepository: repo,
	}
}
