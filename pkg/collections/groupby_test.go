package collections

import (
	"reflect"
	"testing"
)

type mockStruct struct {
	value  int
	key    string
	result int
}

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name string
		in   []interface{}
		f    func(interface{}) (interface{}, interface{})
		want map[interface{}][]interface{}
	}{
		{
			name: "group integers by even and odd",
			in:   []interface{}{1, 2, 3, 4},
			f: func(v interface{}) (interface{}, interface{}) {
				i := v.(int)
				if i%2 == 0 {
					return "even", i
				}
				return "odd", i
			},
			want: map[interface{}][]interface{}{"odd": {1, 3}, "even": {2, 4}},
		},
		{
			name: "group strings by their first letter",
			in:   []interface{}{"apple", "banana", "apricot", "cherry"},
			f: func(v interface{}) (interface{}, interface{}) {
				s := v.(string)
				return rune(s[0]), s
			},
			want: map[interface{}][]interface{}{
				'a': {"apple", "apricot"},
				'b': {"banana"},
				'c': {"cherry"},
			},
		},
		{
			name: "group by nil key",
			in:   []interface{}{"one", "two", nil},
			f: func(v interface{}) (interface{}, interface{}) {
				if v == nil {
					return nil, "nil value"
				}
				return "not nil", v
			},
			want: map[interface{}][]interface{}{
				"not nil": {"one", "two"},
				nil:       {"nil value"},
			},
		},
		{
			name: "empty input slice",
			in:   []interface{}{},
			f: func(v interface{}) (interface{}, interface{}) {
				return v, v
			},
			want: map[interface{}][]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupBy(tt.in, tt.f)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
