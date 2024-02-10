package shell

import (
	"testing"
	"time"
)

func TestBash(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		wantErr bool
	}{
		{"", []string{"ping", "-c", "4", "8.8.8.8"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout, gotStderr, err := Bash(tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bash() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("Bash() gotStdout = %v", gotStdout)
			t.Logf("Bash() gotStderr = %v", gotStderr)
		})
	}
}

func TestBashWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		command []string
		wantErr bool
	}{
		{"", time.Second * 30, []string{"ping", "-c", "4", "8.8.8.8"}, false},
		{"", time.Millisecond, []string{"ping", "-c", "4", "8.8.8.8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout, gotStderr, err := BashWithTimeout(tt.timeout, tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BashWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("BashWithTimeout() gotStdout = %v", gotStdout)
			t.Logf("BashWithTimeout() gotStderr = %v", gotStderr)
		})
	}
}
