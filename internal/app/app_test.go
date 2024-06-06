package app

import (
	"fmt"
	"testing"
)

func Test_getMessage(t *testing.T) {
	type args struct {
		value float64
		prev  float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "value > high threshold",
			args: args{value: highThreshold + 1.2},
			want: fmt.Sprintf(highText, highThreshold+1.2),
		},
		{
			name: "value < high threshold",
			args: args{value: lowThreshold - 1.2},
			want: fmt.Sprintf(lowText, lowThreshold-1.2),
		},
		{
			name: "value = high threshold",
			args: args{value: highThreshold},
			want: "",
		},
		{
			name: "value = low threshold",
			args: args{value: lowThreshold},
			want: "",
		},
		{
			name: "prev value is ok, but value > high threshold",
			args: args{
				value: highThreshold + 1.2,
				prev:  highThreshold,
			},
			want: fmt.Sprintf(highText, highThreshold+1.2),
		},
		{
			name: "prev value is ok, but value < low threshold",
			args: args{
				value: lowThreshold - 1.2,
				prev:  lowThreshold,
			},
			want: fmt.Sprintf(lowText, lowThreshold-1.2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.prev != 0 {
				prevValue = tt.args.prev
			}
			if got := getMessage(tt.args.value); got != tt.want {
				t.Errorf("getMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid date",
			args: args{
				date: "2024-06-05T22:30:37.265Z",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidDate(tt.args.date); got != tt.want {
				t.Errorf("isValidDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
