package slice

import "testing"

func TestEqual(t *testing.T) {
	type args[T comparable] struct {
		s1 []T
		s2 []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{
				s1: []int{},
				s2: []int{},
			},
			want: true,
		},
		{
			name: "nils",
			args: args[int]{
				s1: nil,
				s2: nil,
			},
			want: true,
		},
		{
			name: "mixed_empty",
			args: args[int]{
				s1: []int{},
				s2: nil,
			},
			want: true,
		},
		{
			name: "mixed_empty2",
			args: args[int]{
				s1: nil,
				s2: []int{},
			},
			want: true,
		},
		{
			name: "one_ok",
			args: args[int]{
				s1: []int{27},
				s2: []int{27},
			},
			want: true,
		},
		{
			name: "one_bad",
			args: args[int]{
				s1: []int{27},
				s2: []int{28},
			},
			want: false,
		},
		{
			name: "three_ok",
			args: args[int]{
				s1: []int{27, 1, -3},
				s2: []int{27, 1, -3},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	type args[T comparable] struct {
		s      []T
		prefix []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{
				s:      []int{},
				prefix: []int{},
			},
			want: true,
		},
		{
			name: "nils",
			args: args[int]{
				s:      nil,
				prefix: nil,
			},
			want: true,
		},
		{
			name: "mixed_empty",
			args: args[int]{
				s:      []int{},
				prefix: nil,
			},
			want: true,
		},
		{
			name: "mixed_empty2",
			args: args[int]{
				s:      nil,
				prefix: []int{},
			},
			want: true,
		},
		{
			name: "one_ok",
			args: args[int]{
				s:      []int{27},
				prefix: []int{27},
			},
			want: true,
		},
		{
			name: "one_bad",
			args: args[int]{
				s:      []int{27},
				prefix: []int{28},
			},
			want: false,
		},
		{
			name: "three_ok1",
			args: args[int]{
				s:      []int{27, 1, -3},
				prefix: []int{27},
			},
			want: true,
		},
		{
			name: "three_ok2",
			args: args[int]{
				s:      []int{27, 1, -3},
				prefix: []int{27, 1},
			},
			want: true,
		},
		{
			name: "three_ok3",
			args: args[int]{
				s:      []int{27, 1, -3},
				prefix: []int{27, 1, -3},
			},
			want: true,
		},
		{
			name: "prefix_too_long",
			args: args[int]{
				s:      []int{27, 1, -3},
				prefix: []int{27, 1, 4, 5},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPrefix(tt.args.s, tt.args.prefix); got != tt.want {
				t.Errorf("HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	type args[T comparable] struct {
		s      []T
		suffix []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{
				s:      []int{},
				suffix: []int{},
			},
			want: true,
		},
		{
			name: "nils",
			args: args[int]{
				s:      nil,
				suffix: nil,
			},
			want: true,
		},
		{
			name: "mixed_empty",
			args: args[int]{
				s:      []int{},
				suffix: nil,
			},
			want: true,
		},
		{
			name: "mixed_empty2",
			args: args[int]{
				s:      nil,
				suffix: []int{},
			},
			want: true,
		},
		{
			name: "one_ok",
			args: args[int]{
				s:      []int{27},
				suffix: []int{27},
			},
			want: true,
		},
		{
			name: "one_bad",
			args: args[int]{
				s:      []int{27},
				suffix: []int{28},
			},
			want: false,
		},
		{
			name: "three_ok1",
			args: args[int]{
				s:      []int{27, 1, -3},
				suffix: []int{-3},
			},
			want: true,
		},
		{
			name: "three_ok2",
			args: args[int]{
				s:      []int{27, 1, -3},
				suffix: []int{1, -3},
			},
			want: true,
		},
		{
			name: "three_ok3",
			args: args[int]{
				s:      []int{27, 1, -3},
				suffix: []int{27, 1, -3},
			},
			want: true,
		},
		{
			name: "suffix_too_long",
			args: args[int]{
				s:      []int{27, 1, -3},
				suffix: []int{27, 1, 4, 5},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.args.s, tt.args.suffix); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
