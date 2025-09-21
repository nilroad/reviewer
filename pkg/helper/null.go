package helper

// PtrOrNil returns a pointer to the value if it's not the zero value of its type.
// Otherwise, it returns nil.
func PtrOrNil[T comparable](val T) *T {
	var zero T
	if val == zero {
		return nil
	}

	return &val
}

func ValueOrZero[T comparable](val *T) T {
	var zero T
	if val == nil {
		return zero
	}

	return *val
}
