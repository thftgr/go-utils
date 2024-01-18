package shell

import (
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		wantErr bool
	}{
		{"", []string{"ping", "8.8.8.8"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout, gotStderr, err := Cmd(tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("Cmd() gotStdout = %v", gotStdout)
			t.Logf("Cmd() gotStderr = %v", gotStderr)
		})
	}
}

func TestCmdWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		command []string
		wantErr bool
	}{
		{"", time.Second * 30, []string{"ping", "8.8.8.8"}, false},
		{"", time.Millisecond, []string{"ping", "8.8.8.8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout, gotStderr, err := CmdWithTimeout(tt.timeout, tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CmdWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("CmdWithTimeout() gotStdout = %v", gotStdout)
			t.Logf("CmdWithTimeout() gotStderr = %v", gotStderr)
		})
	}
}

func TestPowershell(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		wantErr bool
	}{
		{"", []string{"ping 8.8.8.8"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout, gotStderr, err := Powershell(tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Powershell() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("Powershell() gotStdout = %v", gotStdout)
			t.Logf("Powershell() gotStderr = %v", gotStderr)
		})
	}
}

func TestPowershellWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		command []string
		wantErr bool
	}{
		{"", time.Second * 30, []string{"ping 8.8.8.8"}, false},
		{"", time.Millisecond, []string{"ping 8.8.8.8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout, gotStderr, err := PowershellWithTimeout(tt.timeout, tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PowershellWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("PowershellWithTimeout() gotStdout = %v", gotStdout)
			t.Logf("PowershellWithTimeout() gotStderr = %v", gotStderr)
		})
	}
}
