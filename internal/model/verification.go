package model

type Verification struct {
	Code Code
}

func NewVerification(codeVal string) (Verification, error) {
	code, err := NewCode(codeVal)
	if err != nil {
		return Verification{}, err
	}

	return Verification{
		Code: code,
	}, nil
}
