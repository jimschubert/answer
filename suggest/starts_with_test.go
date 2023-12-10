package suggest

import (
	"reflect"
	"testing"

	_ "github.com/charmbracelet/x/exp/teatest"
)

func TestStartsWith(t *testing.T) {
	sample := []string{"able", "ablest", "ablative", "abba", "about", "batter", "battering", "battery"}
	type args struct {
		data    []string
		options []StartsWithOpt
	}
	tests := []struct {
		name  string
		args  args
		value string
		want  []string
	}{
		{
			name:  "empty when passed empty string",
			args:  args{data: sample},
			value: "",
			want:  []string{},
		},
		{
			name:  "empty when no matches found",
			args:  args{data: sample},
			value: "car",
			want:  []string{},
		},
		{
			name:  "results when filtered short with default options",
			args:  args{data: sample},
			value: "ab",
			want:  []string{},
		},
		{
			name:  "results when filtered short with modified minimum length",
			args:  args{data: sample, options: []StartsWithOpt{StartsWithMin(2)}},
			value: "ab",
			want:  []string{"able", "ablest", "ablative", "abba", "about"},
		},
		{
			name:  "results when filtered long",
			args:  args{data: sample},
			value: "batter",
			want:  []string{"batter", "battering", "battery"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartsWith(tt.args.data, tt.args.options...)(tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
