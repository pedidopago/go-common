package slice

import (
	"reflect"
	"testing"
)

func TestMergeSlices(t *testing.T) {
	testCases := []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{
			name:     "both slices are empty",
			slice1:   []int{},
			slice2:   []int{},
			expected: []int{},
		},
		{
			name:     "one slice is empty",
			slice1:   []int{1, 2, 3},
			slice2:   []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "no duplicates",
			slice1:   []int{1, 2, 3},
			slice2:   []int{4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "duplicates in slice1",
			slice1:   []int{1, 2, 3, 4},
			slice2:   []int{3, 4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "duplicates in slice2",
			slice1:   []int{1, 2, 3},
			slice2:   []int{3, 4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "duplicates in both slices",
			slice1:   []int{1, 2, 3, 4},
			slice2:   []int{3, 4, 5, 6, 1},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MergeSlices(tc.slice1, tc.slice2)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
