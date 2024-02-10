package shell

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

func DefaultShell(command ...string) (stdout, stderr string, err error) {
	return Bash(command...)
}

func Bash(command ...string) (stdout, stderr string, err error) {
	return run(exec.Command("bash", "-c", strings.Join(command, " ")))
}

func BashWithTimeout(timeout time.Duration, command ...string) (stdout, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return run(exec.CommandContext(ctx, "bash", "-c", strings.Join(command, " ")))
}
