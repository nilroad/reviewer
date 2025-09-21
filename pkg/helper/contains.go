package helper

func Contains[T string | ~int | ~float64](elem T, array []T) bool {
	for _, element := range array {
		if element == elem {
			return true
		}
	}

	return false
}

func SliceToDomain[T interface{ ToDomain() B }, B any](data []T) []B {
	domains := make([]B, 0, len(data))
	for _, v := range data {
		domains = append(domains, v.ToDomain())
	}

	return domains
}
