package model

import (
	"fmt"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type UUID string

func NewUUID() UUID {
	return UUID(uuid.NewString())
}

func ToUUID(id string) (UUID, error) {
	idValueObject := UUID(id)
	if !idValueObject.IsValid() {
		return "", fmt.Errorf("id: %s: %w", id, errs.ErrInvalidArgument)
	}
	return idValueObject, nil
}

func (u UUID) IsValid() bool {
	if _, err := uuid.Parse(u.String()); err != nil {
		return false
	}
	return true
}

func (u UUID) String() string {
	return string(u)
}
