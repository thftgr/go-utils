package shell

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

func DefaultShell(command ...string) (stdout, stderr string, err error) {
	return Cmd(command...)
}

func Cmd(command ...string) (stdout, stderr string, err error) {
	return run(exec.Command("cmd", append([]string{"/c", "chcp", "65001", ">nul", "&&"}, command...)...))
}

func CmdWithTimeout(timeout time.Duration, command ...string) (stdout, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return run(exec.CommandContext(ctx, "cmd", append([]string{"/c", "chcp", "65001", ">nul", "&&"}, command...)...))
}
func Powershell(command ...string) (stdout, stderr string, err error) {
	for i := range command {
		command[i] = "if ($?) { " + command[i] + " }"
	}
	commands := "[void](chcp 65001); " + strings.Join(command, " ")
	cmd := exec.Command("powershell", []string{"-NoProfile", "-Command", commands}...)
	return run(cmd)
}

func PowershellWithTimeout(timeout time.Duration, command ...string) (stdout, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for i := range command {
		command[i] = "if ($?) { " + command[i] + " }"
	}
	commands := "[void](chcp 65001); " + strings.Join(command, " ")
	cmd := exec.CommandContext(ctx, "powershell", []string{"-NoProfile", "-Command", commands}...)
	return run(cmd)
}
