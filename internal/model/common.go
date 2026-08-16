package model

type ValueObject interface {
	String() string
}

func ConvertToValueObjsArray[T any](s []string, f func(string) (T, error)) ([]T, error) {
	valueObjects := make([]T, len(s))
	for i := range s {
		valueObject, err := f(s[i])
		if err != nil {
			return []T{}, err
		}
		valueObjects[i] = valueObject
	}

	return valueObjects, nil
}

func ConvertFromValueObjsArray[T ValueObject](valueObjects []T) []string {
	strs := make([]string, len(valueObjects))
	for i := range valueObjects {
		strs[i] = valueObjects[i].String()
	}

	return strs
}
