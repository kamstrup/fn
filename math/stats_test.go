package fnmath

import (
	"reflect"
	"testing"

	"github.com/kamstrup/fn"
)

func TestStats(t *testing.T) {
	cases := []struct {
		name   string
		input  []int
		expect Stats[int]
	}{
		{
			name:   "empty",
			input:  nil,
			expect: Stats[int]{},
		},
		{
			name:  "zero",
			input: []int{0},
			expect: Stats[int]{
				Count: 1,
			},
		},
		{
			name:  "one",
			input: []int{1},
			expect: Stats[int]{
				Sum:   1,
				Min:   1,
				Max:   1,
				Count: 1,
			},
		}, {
			name:  "minus one",
			input: []int{-1},
			expect: Stats[int]{
				Sum:   -1,
				Min:   -1,
				Max:   -1,
				Count: 1,
			},
		},
		{
			name:  "range",
			input: []int{-1, 0, 1, 2},
			expect: Stats[int]{
				Sum:   2,
				Min:   -1,
				Max:   2,
				Count: 4,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := fn.SliceOf(tc.input)
			opt := fn.Into(Stats[int]{}, MakeStats[int], data)
			if len(tc.input) == 0 {
				if opt.Ok() {
					t.Fatalf("Stats opt should be empty when input is empty")
				}
			} else {
				stats := opt.Must()
				if !reflect.DeepEqual(stats, tc.expect) {
					t.Fatalf("stats mismatch\nexpected: %v\ngot:      %v", tc.expect, stats)
				}
			}
		})
	}

}
