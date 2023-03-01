package slice

func MergeSlices[T comparable](slice1 []T, slice2 []T) []T {
	set := make(map[T]bool)
	merged := []T{}

	for _, value := range slice1 {
		set[value] = true
		merged = append(merged, value)
	}

	for _, value := range slice2 {
		if _, ok := set[value]; !ok {
			set[value] = true
			merged = append(merged, value)
		}
	}

	return merged
}
