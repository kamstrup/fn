package slice

import (
	"reflect"
	"testing"
)

func TestMapping(t *testing.T) {
	even := Mapping([]int{1, 2, 3}, func(i int) int {
		return i * 2
	})
	if !reflect.DeepEqual(even, []int{2, 4, 6}) {
		t.Fatalf("bad results: %v", even)
	}
}

func TestMapIndex(t *testing.T) {
	results := MappingIndex([]int{1, 2, 3}, func(idx, n int) int {
		return n + idx
	})
	if !reflect.DeepEqual(results, []int{1, 3, 5}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestGen(t *testing.T) {
	results := Gen(3, func(idx int) int {
		return idx
	})
	if !reflect.DeepEqual(results, []int{0, 1, 2}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestCopy(t *testing.T) {
	orig := []int{1, 2, 3}
	results := Copy(orig)
	if !reflect.DeepEqual(results, []int{1, 2, 3}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestZero(t *testing.T) {
	orig := []int{1, 2, 3}
	results := Zero(orig)
	if !reflect.DeepEqual(results, []int{0, 0, 0}) || !reflect.DeepEqual(orig, []int{0, 0, 0}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestSortAsc(t *testing.T) {
	data := []int{1, 3, 2}
	SortAsc(data)
	if !reflect.DeepEqual(data, []int{1, 2, 3}) {
		t.Fatalf("bad results: %v", data)
	}
}

func TestSortDesc(t *testing.T) {
	data := []int{1, 3, 2}
	SortDesc(data)
	if !reflect.DeepEqual(data, []int{3, 2, 1}) {
		t.Fatalf("bad results: %v", data)
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		s            []int
		shouldDelete func(int) bool
	}
	type testCase struct {
		name string
		args args
		want []int
	}
	tests := []testCase{
		{
			name: "nil",
			args: args{
				s:            nil,
				shouldDelete: func(n int) bool { return false },
			},
			want: nil,
		},
		{
			name: "empty",
			args: args{
				s:            []int{},
				shouldDelete: func(n int) bool { return false },
			},
			want: []int{},
		},
		{
			name: "del_one",
			args: args{
				s:            []int{1},
				shouldDelete: func(n int) bool { return true },
			},
			want: []int{},
		},
		{
			name: "del_first",
			args: args{
				s:            []int{1, 2, 3},
				shouldDelete: func(n int) bool { return n == 1 },
			},
			want: []int{2, 3},
		},
		{
			name: "del_last",
			args: args{
				s:            []int{1, 2, 3},
				shouldDelete: func(n int) bool { return n == 3 },
			},
			want: []int{1, 2},
		},
		{
			name: "del_mid",
			args: args{
				s:            []int{1, 2, 3},
				shouldDelete: func(n int) bool { return n == 2 },
			},
			want: []int{1, 3},
		},
		{
			name: "del_odd",
			args: args{
				s:            []int{1, 2, 3, 4, 5},
				shouldDelete: func(n int) bool { return n%2 == 1 },
			},
			want: []int{2, 4},
		},
		{
			name: "del_even",
			args: args{
				s:            []int{1, 2, 3, 4, 5},
				shouldDelete: func(n int) bool { return n%2 == 0 },
			},
			want: []int{1, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.s, tt.args.shouldDelete); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type testCase struct {
		name string
		arg  []int
		want []int
	}
	tests := []testCase{
		{
			name: "nil",
			arg:  nil,
			want: nil,
		},
		{
			name: "empty",
			arg:  []int{},
			want: []int{},
		},
		{
			name: "one_start",
			arg:  []int{0, 1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "one_end",
			arg:  []int{1, 2, 3, 0},
			want: []int{1, 2, 3},
		},
		{
			name: "one_both",
			arg:  []int{0, 1, 2, 3, 0},
			want: []int{1, 2, 3},
		},
		{
			name: "many_both_mixed",
			arg:  []int{0, 0, 0, 1, 2, 0, 3, 0, 0, 0, 0, 0},
			want: []int{1, 2, 0, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}
