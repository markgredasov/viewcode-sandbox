package service_verifications

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *VerificationsService) Create(ctx context.Context, verification model.Verification, format *string) (model.Verification, error) {
	if format == nil {
		return model.Verification{}, fmt.Errorf("no format provided: %w", errs.ErrInvalidArgument)
	}

	formatStr := *format
	verification.ID = uuid.NewString()

	switch formatStr {
	case "base64":
		b, err := base64.StdEncoding.DecodeString(verification.Code.String())
		if err != nil {
			return model.Verification{}, fmt.Errorf("decode code value: %w", errs.ErrInvalidArgument)
		}

		decoded, err := model.NewCode(string(b))
		if err != nil {
			return model.Verification{}, err
		}

		verification.Code = decoded

		return s.verificationsRepository.Create(ctx, verification)
	default:
		return model.Verification{}, fmt.Errorf("invalid format: %w", errs.ErrInvalidArgument)
	}
}
