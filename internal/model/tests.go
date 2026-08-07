package model

type Tests struct {
	Code Code
}

func NewTests(codeVal string) (Tests, error) {
	code, err := NewCode(codeVal)
	if err != nil {
		return Tests{}, err
	}

	return Tests{
		Code: code,
	}, nil
}
