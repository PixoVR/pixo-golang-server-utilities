package common

func GetPointer[T any](input T) *T {
	return &input
}

func GetPointers[T any](input []T) []*T {
	pointers := make([]*T, len(input))
	for i, v := range input {
		pointers[i] = &v
	}
	return pointers
}

func GetStructs[T any](input []*T) []T {
	values := make([]T, len(input))
	for i, v := range input {
		if v != nil {
			values[i] = *v
		}
	}
	return values
}

func GetPointerOrDefault[T any](input *T, defaultValue T) *T {
	if input == nil {
		return &defaultValue
	}
	return input
}

func GetValueOrDefault[T any](input *T, defaultValue T) T {
	if input == nil {
		return defaultValue
	}
	return *input
}

func PointerValuesAreEqual[T comparable](value1, value2 *T) bool {
	if value1 == nil || value2 == nil {
		return value1 == value2
	}

	return *value1 == *value2
}
