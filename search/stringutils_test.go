package search

import (
	"testing"
)

const darkness = `Hello darkness, my old friend
I've come to talk with you again
Because a vision softly creeping
Left its seeds while I was sleeping
And the vision that was planted in my brain
Still remains
Within the sound of silence`

func Test_max(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Left is greater", args{100, 10}, 100},
		{"Equal", args{100, 100}, 100},
		{"Right is greater", args{10, 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Test_getLineIndexes(t *testing.T) {

	tests := []struct {
		name, txt string
		want      []int
	}{
		{
			"One line",
			"Hello World :D",
			[]int{0},
		},

		{
			"Empty.",
			"",
			[]int{}, // No lines: empty slice
		},

		{
			"Two Lines",
			"Hello\nWorld :D",
			[]int{0, 6},
		},

		{
			"Empty line at EOF",
			"Hello\nWorld :D\n",
			[]int{0, 6},
		},

		{
			"Windows End of Line (CRLF)",
			"Hello World\r\n:D",
			[]int{0, 13},
		},

		{
			"Windows (CRLF) empty line at EOF",
			"Hello World\r\nFOO\r\n",
			[]int{0, 13},
		},

		{
			"Darkness",
			darkness,
			[]int{0, 30, 63, 96, 132, 176, 190},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLineIndexes(tt.txt); !sliceEq(got, tt.want) {
				t.Errorf("getLineIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bound_convert(t *testing.T) {

	idxs := []int{0, 30, 63, 96, 132, 176, 190}

	tests := []struct {
		name string
		b    bound
		want bound
	}{
		{"Test 1", bound{0, 1}, bound{0, 30}},
		{"Test 2", bound{1, 3}, bound{30, 96}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make the conversion
			tt.b.convert(idxs)
			if tt.b.start != tt.want.start || tt.b.end != tt.want.end {
				t.Errorf("convert() = %v, want %v", tt.b, tt.want)
			}
		})
	}
}
