package shell

import (
	"os/exec"
	"strings"
)

func DefaultShell(command ...string) (stdout, stderr string, err error) {
	return Bash(command...)
}

func Bash(command ...string) (stdout, stderr string, err error) {
	return run(exec.Command("bash", "-c", strings.Join(command, " ")))
}
