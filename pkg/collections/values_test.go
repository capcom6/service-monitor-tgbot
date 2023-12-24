package collections

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestValues(t *testing.T) {
	tests := []struct {
		name string
		in   map[interface{}]interface{}
		want []interface{}
	}{
		{
			name: "non-empty map int to string",
			in:   map[interface{}]interface{}{1: "one", 2: "two", 3: "three"},
			want: []interface{}{"one", "two", "three"}, // Note: order is not guaranteed
		},
		{
			name: "non-empty map string to int",
			in:   map[interface{}]interface{}{"one": 1, "two": 2, "three": 3},
			want: []interface{}{1, 2, 3}, // Note: order is not guaranteed
		},
		{
			name: "empty map",
			in:   map[interface{}]interface{}{},
			want: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Values(tt.in)

			if !reflect.DeepEqual(sortAndReturn(got), sortAndReturn(tt.want)) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to sort the slices since map iteration order is not guaranteed.
func sortAndReturn(slice []interface{}) []interface{} {
	sort.Slice(slice, func(i, j int) bool {
		return fmt.Sprint(slice[i]) < fmt.Sprint(slice[j])
	})
	return slice
}
