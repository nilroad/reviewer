package convertor

import (
	"errors"
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func ConvertToSliceUint64[T Number](n []T) []uint64 {
	result := make([]uint64, len(n))
	for i, num := range n {
		result[i] = uint64(num)
	}

	return result
}

func ConvertToSliceInterface[T Number](n []T) []any {
	result := make([]any, len(n))
	for i, num := range n {
		result[i] = num
	}

	return result
}

func SafeConvertIntToUint64[T interface {
	int | int8 | int16 | int32 | int64
}](i T) (uint64, error) {
	if i < 0 {
		return 0, errors.New("cannot convert negative int64 to uint64")
	}

	return uint64(i), nil
}

func SafeConvertUintToInt[T uint | uint64](i T) (int, error) {
	if i >= math.MaxInt64 {
		return 0, errors.New("cannot convert negative int64 to int")
	}

	return int(i), nil
}
