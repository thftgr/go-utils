package utils

import (
	"strings"
	"testing"
)

func TestGetSourceLine(t *testing.T) {
	tests := []struct {
		name string
		skip int
		want string
	}{
		{
			name: "Zero skip",
			skip: 0,
			want: "utils/runtime_test.go",
		},
		{
			name: "One skip",
			skip: 1,
			want: "testing/testing.go",
		},
		{
			name: "Big skip",
			skip: 1000,
			want: "???",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ans := GetSourceLine(tt.skip)
			file := strings.Split(ans, ":")[0]

			if file != tt.want {
				t.Errorf("got %v; want %v", file, tt.want)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	tests := []struct {
		name string
		skip int
		want string
	}{
		{
			name: "Skip: 0",
			skip: 0,
			want: "runtime_test.go",
		},
		{
			name: "Skip: 9999",
			skip: 9999,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileName(tt.skip); got != tt.want {
				t.Errorf("GetFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
