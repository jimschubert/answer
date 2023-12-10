package suggest

import (
	"reflect"
	"testing"

	_ "github.com/charmbracelet/x/exp/teatest"
)

func TestCalculateLevenshteinDistance(t *testing.T) {
	cases := []struct {
		name string
		str1 string
		str2 string
		want int
	}{
		{
			name: "empty strings",
			str1: "",
			str2: "",
			want: 0,
		},
		{
			name: "first string empty",
			str1: "",
			str2: "hello",
			want: 5,
		},
		{
			name: "second string empty",
			str1: "hello",
			str2: "",
			want: 5,
		},
		{
			name: "both strings identical",
			str1: "hello",
			str2: "hello",
			want: 0,
		},
		{
			name: "both strings different",
			str1: "kitten",
			str2: "sitting",
			want: 3,
		},
		{
			name: "both strings different ignoring case",
			str1: "kiTTen",
			str2: "sitting",
			want: 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := calculateLevenshteinDistance(tc.str1, tc.str2, true); got != tc.want {
				t.Errorf("calculateLevenshteinDistance(%v, %v) = %v; want %v", tc.str1, tc.str2, got, tc.want)
			}
		})
	}
}

func TestLevenshteinDistance(t *testing.T) {
	sample := []string{"kitten", "pumpkin", "sitting", "MELLOW", "world", "yellow"}

	type args struct {
		data    []string
		options []LevenshteinDistanceOpt
	}
	tests := []struct {
		name  string
		args  args
		value string
		want  []string
	}{
		{
			name:  "empty strings return no results",
			args:  args{data: sample},
			value: "",
			want:  []string{},
		},
		{
			name:  "supplies results based on edit distance",
			args:  args{data: sample},
			value: "sitting",
			want:  []string{"kitten", "sitting"},
		},
		{
			name: "supplies results based on edit distance with maximum edit",
			args: args{data: sample, options: []LevenshteinDistanceOpt{
				LevenshteinDistanceMax(2),
			}},
			value: "sitting",
			want:  []string{"sitting"},
		},
		{
			name: "supplies results based on edit distance with minimum edit",
			args: args{data: sample, options: []LevenshteinDistanceOpt{
				LevenshteinDistanceMin(1),
			}},
			value: "sitting",
			want:  []string{"kitten"},
		},
		{
			name: "supplies results based on edit distance while honoring case",
			args: args{data: sample, options: []LevenshteinDistanceOpt{
				LevenshteinDistanceIgnoreCase(false),
			}},
			value: "MELLOW",
			want:  []string{"MELLOW"},
		},

		{
			name: "supplies results based on edit distance with increased edit distance",
			args: args{data: sample, options: []LevenshteinDistanceOpt{
				LevenshteinDistanceMax(6),
			}},
			value: "mellow",
			want:  []string{"kitten", "MELLOW", "world", "yellow"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevenshteinDistance(tt.args.data, tt.args.options...)(tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevenshteinDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}
