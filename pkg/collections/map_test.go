package collections

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	// Define a custom error for testing error handling.
	errTest := errors.New("test error")

	// Test cases defined in a table
	tests := []struct {
		name     string
		input    []int
		function func(int) (string, error)
		expected []string
		err      error
	}{
		{
			name:     "int to string conversion",
			input:    []int{1, 2, 3},
			function: func(i int) (string, error) { return strconv.Itoa(i), nil },
			expected: []string{"1", "2", "3"},
			err:      nil,
		},
		{
			name:  "error on negative numbers",
			input: []int{-1, 2, -3},
			function: func(i int) (string, error) {
				if i < 0 {
					return "", errTest
				}
				return strconv.Itoa(i), nil
			},
			expected: nil,
			err:      errTest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Map(tc.input, tc.function)
			if !reflect.DeepEqual(result, tc.expected) || !errors.Is(err, tc.err) {
				t.Errorf("Test %s - expected result %v, expected error %v, got result %v, got error %v", tc.name, tc.expected, tc.err, result, err)
			}
		})
	}
}
